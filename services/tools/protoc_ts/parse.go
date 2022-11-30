package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// todo: optimize: use Finite State Machine
func (g *generator) parse(fileName string) *pbFile {
	file, err := os.Open(inputDir + fileName)
	if err != nil {
		log.Fatalln("open file failed, error: ", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatalln("close file failed, error: ", err)
		}
	}()

	pbFileIns := &pbFile{
		name:         strings.TrimSuffix(fileName, protoFileSuffix),
		payloadStack: &payloadStack{},
	}

	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		content := fileScanner.Text()

		reMatchedFlag, matchedData := g.reMatched(content)

		switch reMatchedFlag {
		case REMatched_UnMatched:
			continue
		case REMatched_Import:
			pbFileIns.imports = append(pbFileIns.imports, &dependence{from: matchedData[1]})
		case REMatched_MessageStart:
			pbFileIns.payloadStack.checkNestedLevel()

			pbFileIns.payloadStack.push(&payload{
				name: matchedData[1],
				typ:  PayloadType_Class,
			})
		case REMatched_MessageField:
			fieldType := matchedData[2]
			if strings.Contains(fieldType, ".") {
				// handle import structure
				fieldTypeSplit := strings.Split(fieldType, ".") // 0: package name 1: structure
				pbFileIns.importStructure(fieldTypeSplit[0], fieldTypeSplit[1])
			}

			pbFileIns.payloadStack.appendField(&field{
				name:    matchedData[3],
				typ:     fieldType,
				isArray: len(matchedData[1]) > 0,
			})
		case REMatched_EnumStart:
			pbFileIns.payloadStack.checkNestedLevel()

			pbFileIns.payloadStack.push(&payload{
				name: matchedData[1],
				typ:  PayloadType_Enum,
			})
		case REMatched_EnumField:
			pbFileIns.payloadStack.appendField(&field{
				name: matchedData[1],
				typ:  PayloadType_Enum,
			})
		case REMatched_End:
			if pbFileIns.payloadStack.num <= 0 { // skip end of 'service'
				break
			}

			finishedPayload := pbFileIns.payloadStack.pop()
			if pbFileIns.payloadStack.num <= 0 { // end of a first-level 'message' or 'enum'
				pbFileIns.payloads = append(pbFileIns.payloads, finishedPayload)
			} else {
				pbFileIns.payloadStack.appendPayload(finishedPayload)
			}
		}
	}
	if err != nil {
		log.Fatalln("parse file failed, error: ", err)
	}

	return pbFileIns
}

func (g *generator) reMatched(content string) (REMatched, []string) {
	reMatchedFlag := REMatched_UnMatched
	matchedData := make([]string, 0)

	if matchedData = g.reEngineIns.importREEngine.FindStringSubmatch(content); len(matchedData) > 0 {
		reMatchedFlag = REMatched_Import
	} else if matchedData = g.reEngineIns.messageStartREEngine.FindStringSubmatch(content); len(matchedData) > 0 {
		reMatchedFlag = REMatched_MessageStart
	} else if matchedData = g.reEngineIns.messageFieldREEngine.FindStringSubmatch(content); len(matchedData) > 0 {
		reMatchedFlag = REMatched_MessageField
	} else if matchedData = g.reEngineIns.enumStartREEngine.FindStringSubmatch(content); len(matchedData) > 0 {
		reMatchedFlag = REMatched_EnumStart
	} else if matchedData = g.reEngineIns.enumFieldREEngine.FindStringSubmatch(content); len(matchedData) > 0 {
		reMatchedFlag = REMatched_EnumField
	} else if matchedData = g.reEngineIns.endREEngine.FindStringSubmatch(content); len(matchedData) > 0 {
		reMatchedFlag = REMatched_End
	}

	return reMatchedFlag, matchedData
}
