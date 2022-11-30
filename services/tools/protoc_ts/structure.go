package main

import (
	"regexp"
)

type pbFile struct {
	name     string
	imports  []*dependence
	payloads []*payload

	payloadStack *payloadStack
}

type dependence struct {
	from       string   // e.g. "common", 'common.proto' remove extension name
	structures []string // e.g. []string{"Error","Pagination"}
}

type payload struct {
	name        string // e.g. "Note", from 'message Note {'
	typ         string // e.g. "namespace"/"class"/"enum"
	indentation int    // payload indentation, field indentation is 'v'+2
	fields      []*field
	children    []*payload
}

// field include 'message' and 'enum' field
type field struct {
	name    string // e.g. "note_id" from 'string note_id = 1;'
	typ     string // e.g. "string"/"ListRule"/"enum"
	isArray bool   // e.g. "repeated Data data = 1;"
}

type reEngines struct {
	importREEngine       *regexp.Regexp
	messageStartREEngine *regexp.Regexp
	messageFieldREEngine *regexp.Regexp // make sure one content can match both 'message field' and 'enum field', or use status
	enumStartREEngine    *regexp.Regexp
	enumFieldREEngine    *regexp.Regexp // make sure one content can match both 'message field' and 'enum field', or use status
	endREEngine          *regexp.Regexp // end of 'message' and 'enum'
}

// importStructure append a 'structure' to an exist 'from'
func (pb *pbFile) importStructure(from string, structure string) {
	for i := range pb.imports {
		if pb.imports[i].from != from {
			continue
		}

		isExist := false
		for j := range pb.imports[i].structures {
			if pb.imports[i].structures[j] != structure {
				continue
			}

			isExist = true
			break
		}

		if !isExist {
			pb.imports[i].structures = append(pb.imports[i].structures, structure)
		}
	}
}

// payloadStack the zero stack is empty and ready for use
type payloadStack struct {
	data []*payload

	// num of array, also used as nested level
	num int
}

func (s *payloadStack) push(data *payload) {
	if s.num >= len(s.data) { // need increase data array
		s.data = append(s.data, &payload{})
	}

	if data.indentation < 1 {
		data.indentation = s.num * indentationSpace
	}

	s.data[s.num] = data
	s.num++
}

func (s *payloadStack) pop() *payload {
	if s.num < 1 {
		return nil
	}

	s.num--

	return s.data[s.num]
}

// appendField append field to last payload
func (s *payloadStack) appendField(data *field) {
	s.data[s.num-1].fields = append(s.data[s.num-1].fields, data)
}

// appendPayload append payload to last payload
func (s *payloadStack) appendPayload(data *payload) {
	s.data[s.num-1].children = append(s.data[s.num-1].children, data)
}

// checkNestedLevel check payload nested level, when match a new 'start', if it is not first-level start,
// upgrade its upper payload type to 'namespace'
func (s *payloadStack) checkNestedLevel() {
	if s.num > 0 {
		s.data[s.num-1].typ = PayloadType_Namespace
	}
}
