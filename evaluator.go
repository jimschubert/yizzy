package yizzy

import (
	"bufio"
	"container/list"
	"io"
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/jimschubert/yizzy/models"
)

type MigrationEvaluator struct {
	treeNavigator yqlib.DataTreeNavigator
	treeCreator   yqlib.ExpressionParser
	migration     *models.MigrationDoc
	envContext    *map[string]string
	currentOps    *models.MigrationEntry
}

func NewMigrationEvaluator(doc *models.MigrationDoc) yqlib.Evaluator {
	return &MigrationEvaluator{treeNavigator: yqlib.NewDataTreeNavigator(), treeCreator: yqlib.NewExpressionParser(), migration: doc}
}

func (y *MigrationEvaluator) EvaluateFiles(expression string, filenames []string, printer yqlib.Printer) error {
	fileIndex := 0

	var allDocuments = list.New()
	for _, filename := range filenames {
		var reader io.Reader
		if filename == "-" {
			reader = bufio.NewReader(os.Stdin)
		} else {
			if r, err := os.Open(filename); err != nil {
				return err
			} else {
				reader = r
			}
		}
		fileDocuments, err := readDocuments(reader, filename, fileIndex)
		if err != nil {
			return err
		}
		allDocuments.PushBackList(fileDocuments)
		fileIndex = fileIndex + 1
	}

	m := make(map[string]string)

	if y.migration.Env != nil {
		for varName, expr := range *y.migration.Env {
			nodes, err := y.EvaluateCandidateNodes(expr, allDocuments)
			if err != nil {
				return err
			}

			if nodes != nil && nodes.Len() > 0 {
				node := nodes.Front().Value.(*yqlib.CandidateNode)
				if node != nil {
					m[varName] = node.Node.Value
				}
			}
		}
	}

	y.envContext = &m

	// TODO: Loop all operations and apply actions against the target expressions. Finally, query one last time against the modified document.

	for _, ops := range *y.migration.Operations {
		y.currentOps = &ops
		var selector = "."
		if ops.Selector != nil {
			selector = *ops.Selector
		}
		_, err := y.EvaluateCandidateNodes(selector, allDocuments)
		if err != nil {
			return err
		}
	}

	matches, err := y.EvaluateCandidateNodes(expression, allDocuments)
	if err != nil {
		return err
	}

	return printer.PrintResults(matches)
}

func (y *MigrationEvaluator) processOperations(matchingNodes *list.List) error {
	return withEnvs(*y.envContext, func() error {
		return nil
	})
}

func withEnvs(envs map[string]string, process func() error) error {
	var originalEnvs = make(map[string]string)
	for key, value := range envs {
		if key != "" {
			originalEnvs[key] = value
			err := os.Setenv(key, value)
			if err != nil {
				return err
			}
		}
	}

	defer func() {
		if len(originalEnvs) > 0 {
			for key, value := range originalEnvs {
				_ = os.Setenv(key, value)
			}
		}
	}()

	return process()
}

func (y *MigrationEvaluator) EvaluateNodes(expression string, nodes ...*yaml.Node) (*list.List, error) {
	inputCandidates := list.New()
	for _, node := range nodes {
		inputCandidates.PushBack(&yqlib.CandidateNode{Node: node})
	}
	return y.EvaluateCandidateNodes(expression, inputCandidates)
}

func (y *MigrationEvaluator) EvaluateCandidateNodes(expression string, inputCandidateNodes *list.List) (*list.List, error) {
	var matchingNodes *list.List = nil
	if expressionNode, err := y.treeCreator.ParseExpression(expression); err == nil {
		var envContext = make(map[string]string)
		if y.envContext != nil {
			for key, value := range *y.envContext {
				envContext[key] = value
			}
		}
		err = withEnvs(envContext, func() error {
			context, err := y.treeNavigator.GetMatchingNodes(yqlib.Context{MatchingNodes: inputCandidateNodes}, expressionNode)
			if err != nil {
				return err
			}
			matchingNodes = context.MatchingNodes

			// I acknowledge this is weird, but yq isn't open for extension at this point and I'd otherwise need to copy much of the lib
			if y.currentOps != nil && matchingNodes != nil && matchingNodes.Len() > 0 {

				var value = ""
				ops := y.currentOps
				if ops != nil {
					// Eval supports query operators such as del, add, etc.
					// This works by evaluating the expression within the context of the matched nodes.
					if ops.Eval != nil && *ops.Eval != "" {
						evalExpression, err := y.treeCreator.ParseExpression(*ops.Eval)
						if err != nil {
							return err
						}
						context, err = y.treeNavigator.GetMatchingNodes(yqlib.Context{MatchingNodes: matchingNodes}, evalExpression)
						if err != nil {
							return err
						}
						applyScalarToNode := func(node *yaml.Node) {
							for e := matchingNodes.Front(); e != nil; e = e.Next() {
								selected := e.Value.(*yqlib.CandidateNode)
								if selected.Path != nil {
									selected.Node.Value = node.Value
									selected.Node.Tag = node.Tag
								}
							}
						}
						for e := context.MatchingNodes.Front(); e != nil; e = e.Next() {
							node := e.Value.(*yqlib.CandidateNode)
							if ops.ValueType != nil && node != nil {
								if node.Path != nil || (node.Path == nil && node.Node.Kind == yaml.ScalarNode) {
									node.Node.Tag = *ops.ValueType
								}
							}

							if node != nil && node.Node != nil && node.Node.Kind == yaml.ScalarNode {
								// a calculated field needs to be set back onto the node(s) found by selector
								applyScalarToNode(node.Node)
							}
						}
					} else {
						for e := matchingNodes.Front(); e != nil; e = e.Next() {
							node := e.Value.(*yqlib.CandidateNode)
							value = *ops.Value
							if node != nil && node.Path != nil {
								switch ops.GetOperation() {
								case "delete":
									break
								case "modify":
									if node.Path != nil && node.Node != nil {
										node.Node.Value = value
										if ops.ValueType != nil {
											node.Node.Tag = *ops.ValueType
										}
									}
									break
								}
							}
						}
					}
				}
			}
			return nil
		})

		return matchingNodes, err
	} else {
		return matchingNodes, err
	}
}

func readDocuments(reader io.Reader, filename string, fileIndex int) (*list.List, error) {
	decoder := yaml.NewDecoder(reader)
	inputList := list.New()
	var currentIndex uint = 0

	for {
		var dataBucket yaml.Node
		errorReading := decoder.Decode(&dataBucket)

		if errorReading == io.EOF {
			switch reader := reader.(type) {
			case *os.File:
				err := reader.Close()
				if err != nil {
					log.WithError(err).Error("Failed to close file.")
				}
			}
			return inputList, nil
		} else if errorReading != nil {
			return nil, errorReading
		}
		candidateNode := &yqlib.CandidateNode{
			Document:         currentIndex,
			Filename:         filename,
			Node:             &dataBucket,
			FileIndex:        fileIndex,
			EvaluateTogether: true,
		}

		inputList.PushBack(candidateNode)

		currentIndex = currentIndex + 1
	}
}
