package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"

	"github.com/jimschubert/yizzy"
)

var version = ""
var date = ""
var commit = ""
var projectName = ""

var opts struct {
	InputFile           string `short:"f" long:"file" description:"The file to process" value-name:"FILE"`
	MigrationsDirectory string `short:"d" long:"dir" description:"The directory where migrations reside"`
	InPlace             bool   `long:"in-place" description:"Writes a file in place"`
	Version             bool   `short:"v" long:"version" description:"Display version information"`
}

const parseArgs = flags.HelpFlag | flags.PassDoubleDash

func main() {
	parser := flags.NewParser(&opts, parseArgs)
	_, err := parser.Parse()
	if err != nil {
		flagError := err.(*flags.Error)
		if flagError.Type == flags.ErrHelp {
			parser.WriteHelp(os.Stdout)
			return
		}

		if flagError.Type == flags.ErrUnknownFlag {
			_, _ = fmt.Fprintf(os.Stderr, "%s. Please use --help for available options.\n", strings.Replace(flagError.Message, "unknown", "Unknown", 1))
			return
		}
		_, _ = fmt.Fprintf(os.Stderr, "Error parsing command line options: %s\n", err)
		return
	}

	if opts.Version {
		fmt.Printf("%s %s (%s) %s\n", projectName, version, commit, date)
		return
	}

	initLogging()

	application := yizzy.YamlMigrator{
		InputFile:           opts.InputFile,
		MigrationsDirectory: opts.MigrationsDirectory,
		InPlace:             opts.InPlace,
	}
	err = application.Run()
	if err != nil {
		log.WithError(err).Errorf("execution failed.")
		return
	}
}

func initLogging() {
	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevel = "error"
	}
	ll, err := log.ParseLevel(logLevel)
	if err != nil {
		ll = log.DebugLevel
	}
	log.SetLevel(ll)
	log.SetOutput(os.Stderr)
}
