package main

import (
	"bufio"
	"fmt"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"github.com/pkg/errors"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type pbFile struct {
	dir      string
	name     string
	imports  []*dependence
	payloads []*payload

    reEngineIns *reEngines
}

type dependence struct {
	from       string   // e.g. "common", 'common.proto' remove extension name
	structures []string // e.g. []string{"Error","Pagination"}
}

type payload struct {
	name        string // e.g. "Note", from 'message Note {'
	typ         string // e.g. "namespace"/"export class"/"export enum"
	indentation int    // payload indentation, field indentation is 'v'+2
	fields      []*field
	children    []*payload
}

// field include 'message' and 'enum' field
type field struct {
	name    string // e.g. "note_id" from 'string note_id = 1;'
	typ     string // e.g. "string"/"enum"(for enum number)
	isArray bool   // e.g. "repeated Data data = 1;"
	number  int    // e.g. "1" from 'string note_id = 1;'
}

type reEngines struct {
	importREEngine       *regexp.Regexp
	messageStartREEngine *regexp.Regexp
	messageFieldREEngine *regexp.Regexp // make sure one content can match both 'message field' and 'enum field', or use status
	enumStartREEngine    *regexp.Regexp
	enumFieldREEngine    *regexp.Regexp // make sure one content can match both 'message field' and 'enum field', or use status
	endREEngine          *regexp.Regexp // end of 'message' and 'enum'
}

var reEngineIns *reEngines

func init() {
	// from src/regexp/syntax/doc.go
	// \w: [0-9A-Za-z_]
	// \s: [\t\n\f\r ], mainly use to match any times of ' '
	importREEngine, err := regexp.Compile(`import\s*"(\w+)\.proto"\s*;`)                            // e.g. 'import "common.proto";', get 'common'
	messageStartREEngine, err2 := regexp.Compile(`message\s+(\w+)\s*{`)                             // e.g. 'message xxx {', get 'xxx'
	messageFieldREEngine, err3 := regexp.Compile(`(repeated)?\s*([\w\.]+)\s+(\w+)\s*=\s*(\d+)\s*;`) // e.g. 'common.Error err = 1;', get 'common.Error'/'err'/'1'
	enumStartREEngine, err4 := regexp.Compile(`enum\s+(.*)\s*{`)                                    // e.g. 'enum xxx {', get 'xxx'
	enumFieldREEngine, err5 := regexp.Compile(`(\w+)\s*=\s*(\d+)\s*;`)                              // e.g. 'UNSPECIFIED = 0;', get 'UNSPECIFIED'/'0'
	endREEngine, err6 := regexp.Compile(`}`)
	if err != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
		log.Fatalln("init RE engine failed", utils.ErrorsToString(err, err2, err3, err4, err5, err6))
	}

	reEngineIns = &reEngines{
		importREEngine:       importREEngine,
		messageStartREEngine: messageStartREEngine,
		messageFieldREEngine: messageFieldREEngine,
		enumStartREEngine:    enumStartREEngine,
		enumFieldREEngine:    enumFieldREEngine,
		endREEngine:          endREEngine,
	}
}

// todo: optimize: use Finite State Machine
func parsePBFile(filePath string, fileName string) *pbFile {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "open file failed"))
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatalln(errors.Wrap(err, "close file failed"))
		}
	}()

	pbFileIns := &pbFile{
		name: fileName,
	}
	payloadStack := make([]*payload, 0)
	nestedLevel := 0 // for nested 'message' and 'enum'

	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		content := fileScanner.Text()
		matchedData := make([]string, 0)

		if matchedData = reEngineIns.importREEngine.FindStringSubmatch(content); len(matchedData) > 0 {
			// import
			pbFileIns.imports = append(pbFileIns.imports, &dependence{from: matchedData[1]})
		} else if matchedData = reEngineIns.messageStartREEngine.FindStringSubmatch(content); len(matchedData) > 0 {
			// message start
			if nestedLevel > 0 {
				payloadStack[len(payloadStack)-1].typ = "namespace"
			}

			payloadStack = append(payloadStack, &payload{
				name:        matchedData[1],
				typ:         "class",
				indentation: nestedLevel * 2,
			})

			nestedLevel++
		} else if matchedData = reEngineIns.messageFieldREEngine.FindStringSubmatch(content); len(matchedData) > 0 {
			// message field
			var fieldNum int
			fieldNum, err = strconv.Atoi(matchedData[4])
			if err != nil {
				break
			}

			fieldType := matchedData[2]
			if strings.Contains(fieldType, ".") {
				// handle import structure
				fieldTypeSplit := strings.Split(fieldType, ".") // 0: package name 1: structure
				for i := range pbFileIns.imports {
					if pbFileIns.imports[i].from != fieldTypeSplit[0] {
						continue
					}

					isExist := false
					for j := range pbFileIns.imports[i].structures {
						if pbFileIns.imports[i].structures[j] == fieldTypeSplit[1] {
							isExist = true
							break
						}
					}

					if !isExist {
						pbFileIns.imports[i].structures = append(pbFileIns.imports[i].structures, fieldTypeSplit[1])
					}
				}
			}

			payloadStack[len(payloadStack)-1].fields = append(payloadStack[len(payloadStack)-1].fields, &field{
				name:    matchedData[3],
				typ:     fieldType,
				isArray: len(matchedData[1]) > 0,
				number:  fieldNum,
			})
		} else if matchedData = reEngineIns.enumStartREEngine.FindStringSubmatch(content); len(matchedData) > 0 {
			// enum start
			if nestedLevel > 0 {
				payloadStack[len(payloadStack)-1].typ = "namespace"
			}

			payloadStack = append(payloadStack, &payload{
				name:        matchedData[1],
				typ:         "enum",
				indentation: nestedLevel * 2,
			})

			nestedLevel++
		} else if matchedData = reEngineIns.enumFieldREEngine.FindStringSubmatch(content); len(matchedData) > 0 {
			// enum field
			var fieldNum int
			fieldNum, err = strconv.Atoi(matchedData[2])
			if err != nil {
				break
			}

			payloadStack[len(payloadStack)-1].fields = append(payloadStack[len(payloadStack)-1].fields, &field{
				name:   matchedData[1],
				typ:    "enum",
				number: fieldNum,
			})
		} else if matchedData = reEngineIns.endREEngine.FindStringSubmatch(content); nestedLevel > 0 && len(matchedData) > 0 { // skip '}' of service
			// end of 'message' and 'enum'
			nestedLevel--
			if nestedLevel <= 0 {
				pbFileIns.payloads = append(pbFileIns.payloads, payloadStack[len(payloadStack)-1])
				payloadStack = payloadStack[:len(payloadStack)-1]
			} else {
				finishedPayload := payloadStack[len(payloadStack)-1]
				payloadStack = payloadStack[:len(payloadStack)-1]

				payloadStack[len(payloadStack)-1].children = append(payloadStack[len(payloadStack)-1].children, finishedPayload)
			}
		}
	}
	if err != nil {
		log.Fatalln(errors.Wrap(err, "parse file failed"))
	}

	return pbFileIns
}

func generateTS(pbFileIns *pbFile) {
	var content []byte
	content = append(content, "// Generate File, should not edit.\n// Author: mario.\n\n"...)
	for i := range pbFileIns.imports {
		content = append(content, generateImport(pbFileIns.imports[i])...)
	}
	content = append(content, '\n')

	for i := range pbFileIns.payloads {
		content = append(content, generatePayload(pbFileIns.payloads[i])...)
	}
	content = append(content, '\n')

	err := os.WriteFile(fmt.Sprintf("./ts/%s.pb.ts", pbFileIns.name), content[:len(content)-2], 0777)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "write file failed"))
	}
}

func generateImport(data *dependence) string {
	var structureList []byte
	for i := range data.structures {
		structureList = append(structureList, ", "+data.structures[i]...)
	}

	if len(structureList) > 0 {
		structureList = structureList[2:]
	}

	return fmt.Sprintf("import { %s } from \"./%s.pb\"\n", structureList, data.from)
}

func generatePayload(data *payload) []byte {
	var content []byte

	for i := 0; i < data.indentation; i++ {
		content = append(content, ' ')
	}
	content = append(content, "export "+data.typ+" "+data.name+" {\n"...)

	for i := range data.fields {
		var line []byte
		for j := 0; j < data.indentation+2; j++ {
			line = append(line, ' ')
		}

		defaultValue := ``
		switch data.fields[i].typ {
		case "string":
			defaultValue = `""`
		case "bool":
			defaultValue = `false`
		case "int32", "int64", "uint32", "uint64", "sint32", "sint64":
			defaultValue = `0`
		case "bytes":
			defaultValue = `Blob`
		default: // self-define type
			typ := strings.Split(data.fields[i].typ, ".")
			defaultValue = typ[len(typ)-1]
		}

		if data.fields[i].isArray {
			defaultValue = fmt.Sprintf("Array<%s>", defaultValue)
		}

		lineFormat := ""
		if data.fields[i].typ == "enum" {
			lineFormat = fmt.Sprintf("%s,\n", data.fields[i].name)
		} else {
			lineFormat = fmt.Sprintf("%s: %s;\n", data.fields[i].name, defaultValue)
		}

		line = append(line, lineFormat...)

		content = append(content, line...)
	}

	for i := range data.children {
		content = append(content, generatePayload(data.children[i])...)
	}

	for i := 0; i < data.indentation; i++ {
		content = append(content, ' ')
	}
	content = append(content, "}\n\n"...)

	return content
}
