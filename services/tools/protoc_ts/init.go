package main

import (
	"flag"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	help             bool
	inputDir         string
	outputDir        string
	indentationSpace int
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.StringVar(&inputDir, "d", "", "input dir\n"+
		"we will handle all '*.proto' files in it\n"+
		"default '[current dir]'")
	flag.StringVar(&outputDir, "o", "", "output dir, default '[input dir]/ts'")
	flag.IntVar(&indentationSpace, "i", 2, "indentation of generate file")

	flag.Parse()

	if help {
		log.Println("Options: ")
		flag.PrintDefaults()
		os.Exit(0)
	}

	parseDir()

	generatorIns.reEngineIns = initREEngines()
}

func parseDir() {
	var err error

	if len(inputDir) < 1 {
		// input dir is empty, use default value
		inputDir, err = os.Getwd()
		if err != nil {
			log.Fatalln("get dir failed, error: ", err)
		}
	}

	// test if input dir is valid
	_, err = os.ReadDir(inputDir)
	if err != nil {
		log.Fatalln("invalid input dir, error: ", err)
	}

	inputDir = formatDir(inputDir)

	if len(outputDir) < 1 {
		// output dir is empty, use default value
		outputDir = formatDir(inputDir) + defaultOutputDir
	}

	// generate output dir if not exist
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalln("'mkdir' on output dir failed, error: ", err)
	}

	outputDir = formatDir(outputDir)
}

// formatDir replace windows dir sep ("\\" and "\") to "/", and make sure dir is end with "/"
func formatDir(dir string) string {
	// for windows
	dir = strings.ReplaceAll(dir, "\\\\", "/")
	dir = strings.ReplaceAll(dir, "\\", "/")

	return utils.FormatDirSuffix(dir)
}

func initREEngines() *reEngines {
	// from src/regexp/syntax/doc.go
	// \w: [0-9A-Za-z_]
	// \s: [\t\n\f\r ], mainly use to match any times of ' '
	importREEngine, err := regexp.Compile(importREExp)
	messageStartREEngine, err2 := regexp.Compile(messageStartREExp)
	messageFieldREEngine, err3 := regexp.Compile(messageFieldREExp)
	enumStartREEngine, err4 := regexp.Compile(enumStartREExp)
	enumFieldREEngine, err5 := regexp.Compile(enumFieldREExp)
	endREEngine, err6 := regexp.Compile(endREExp)
	if err != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
		log.Fatalln("init RE engine failed, error: ", utils.ErrorsToString(err, err2, err3, err4, err5, err6))
	}

	return &reEngines{
		importREEngine:       importREEngine,
		messageStartREEngine: messageStartREEngine,
		messageFieldREEngine: messageFieldREEngine,
		enumStartREEngine:    enumStartREEngine,
		enumFieldREEngine:    enumFieldREEngine,
		endREEngine:          endREEngine,
	}
}
