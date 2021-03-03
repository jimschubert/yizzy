package yizzy

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/op/go-logging.v1"

	"github.com/jimschubert/yizzy/models"
)

// YamlMigrator is an application which applies all migrations defined in MigrationsDirectory to InputFile.
type YamlMigrator struct {
	InputFile           string
	MigrationsDirectory string
	InPlace             bool
}

// Run a migration.
func (y *YamlMigrator) Run() error {
	if err := y.validateFile(); err != nil {
		return err
	}
	if err := y.validateMigrationsDirectory(); err != nil {
		return err
	}
	if !y.InPlace {
		return fmt.Errorf("only inplace file writes are currently supported")
	}

	resetYqLogger()

	docs := make([]*models.MigrationDoc, 0)

	migrationFiles, err := os.ReadDir(y.MigrationsDirectory)
	if err != nil {
		return err
	}
	for _, file := range migrationFiles {
		if !file.IsDir() {
			migrationDocument := &models.MigrationDoc{}
			err = migrationDocument.Load(path.Join(y.MigrationsDirectory, file.Name()))
			if err != nil {
				return err
			}

			docs = append(docs, migrationDocument)
		}
	}

	if len(docs) == 0 {
		return fmt.Errorf("no migrations found in %s", y.MigrationsDirectory)
	}

	for _, migration := range docs {
		err := y.processFile(migration)
		if err != nil {
			return err
		}
	}
	return nil
}

func (y *YamlMigrator) processFile(migration *models.MigrationDoc) error {
	writeInPlaceHandler := yqlib.NewWriteInPlaceHandler(y.InputFile)
	outFile, err := writeInPlaceHandler.CreateTempFile()
	if err != nil {
		return err
	}
	defer func() { writeInPlaceHandler.FinishWriteInPlace(true) }()

	printer := yqlib.NewPrinter(outFile, false, false, false, 2, true)
	eval := NewMigrationEvaluator(migration)
	err = eval.EvaluateFiles(`. | ... style= ""`, []string{y.InputFile}, printer)
	return nil
}

func resetYqLogger() {
	var yqLogger = logging.MustGetLogger("yq-lib")
	b := logging.SetBackend(logging.NewLogBackend(os.Stderr, "", log.LstdFlags))
	b.SetLevel(logging.ERROR, "")
	logging.SetFormatter(logging.DefaultFormatter)
	yqLogger.SetBackend(b)
}

func (y *YamlMigrator) validateFile() error {
	if _, err := os.Stat(y.InputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: '%s'", y.InputFile)
	}
	return nil
}

func (y *YamlMigrator) validateMigrationsDirectory() error {
	stat, err := os.Stat(y.MigrationsDirectory);
	if os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: '%s'", y.MigrationsDirectory)
	}
	if !stat.IsDir() {
		return fmt.Errorf("migrations directory appears to be a file: '%s'", y.MigrationsDirectory)
	}
	return nil
}