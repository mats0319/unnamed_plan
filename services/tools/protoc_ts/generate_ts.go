package main

import (
	"fmt"
	"strings"
)

func tsImportTemplate(structureList string, fileName string) string {
	return fmt.Sprintf("import { %s } from \"./%s%s\"\n", structureList, fileName, defaultImportFileName)
}

func tsPayloadDeclareTemplate(typ string, name string) string {
	return fmt.Sprintf("export %s %s {\n", typ, name)
}

func tsFieldTemplate(data *field, payloadList []*payload) string {
	if data.typ == PayloadType_Enum { // 'enum' type
		return fmt.Sprintf("%s,\n", data.name)
	}

	// basic type or self-define type
	tsParam, ok := tsTypeMap[data.typ]
	if !ok { // self-define type
		typSplit := strings.Split(data.typ, ".")
		tsParam = struct {
			typ          string
			defaultValue string
		}{
			typ: typSplit[len(typSplit)-1],
		}

		if isEnum(payloadList, data.typ) { // self-define type is 'enum'
			tsParam.defaultValue = fmt.Sprintf("%s.%s", tsParam.typ, enumDefaultValue)
		} else {
			tsParam.defaultValue = fmt.Sprintf("new %s()", tsParam.typ)
		}
	}
	if data.isArray {
		tsParam.typ = fmt.Sprintf("Array<%s>", tsParam.typ)
		tsParam.defaultValue = fmt.Sprintf("new %s()", tsParam.typ)
	}

	return fmt.Sprintf("%s: %s = %s;\n", data.name, tsParam.typ, tsParam.defaultValue)
}

// tsTypeMap proto-type -> ts type
var tsTypeMap = map[string]struct {
	typ          string
	defaultValue string
}{
	"string": {typ: "string", defaultValue: `""`},
	"bool":   {typ: "boolean", defaultValue: "false"},
	"bytes":  {typ: "Blob", defaultValue: "new Blob()"},
	"int32":  {typ: "number", defaultValue: "0"},
	"int64":  {typ: "number", defaultValue: "0"},
	"uint32": {typ: "number", defaultValue: "0"},
	"uint64": {typ: "number", defaultValue: "0"},
}

// todo: optimize(?): use template
func (pb *pbFile) generateTSContent() []byte {
	content := []byte(generateFileDeclaration)
	content = append(content, pb.generateTSImport()...)
	content = append(content, pb.generateTSPayload()...)

	return content
}

func (pb *pbFile) generateTSImport() (res []byte) {
	for i := range pb.imports {
		var structureList []byte
		for j := range pb.imports[i].structures {
			structureList = append(structureList, ", "+pb.imports[i].structures[j]...)
		}

		if len(structureList) > 0 {
			structureList = structureList[2:]
		}

		res = append(res, tsImportTemplate(string(structureList), pb.imports[i].from)...)
	}

	if len(res) > 0 {
		res = append(res, '\n')
	}

	return res
}

func (pb *pbFile) generateTSPayload() (res []byte) {
	for i := range pb.payloads {
		res = append(res, pb.generatePayloadItem(pb.payloads[i])...)
	}

	if len(res) > 1 { // del suffix '\n' of last first-level payload
		res = res[:len(res)-1]
	}

	return
}

func (pb *pbFile) generatePayloadItem(data *payload) []byte {
	content := whiteSpaces(data.indentation)
	content = append(content, tsPayloadDeclareTemplate(data.typ, data.name)...)

	for i := range data.fields {
		line := whiteSpaces(data.indentation + indentationSpace)
		line = append(line, tsFieldTemplate(data.fields[i], pb.payloads)...)

		content = append(content, line...)
	}

	for i := range data.children {
		content = append(content, pb.generatePayloadItem(data.children[i])...)
	}

	content = append(content, whiteSpaces(data.indentation)...)
	content = append(content, "}\n\n"...)

	return content
}
