package main

// some default value, normally not modify
const (
	defaultOutputDir           = "ts"
	defaultOutputFileExtension = ".pb.ts"
	defaultImportFileName      = ".pb" // ignore '.ts' in ts files

	generateFileDeclaration = "// Generate File, Should not Edit.\n// Author: mario.\n\n"

	protoFileSuffix  = ".proto"
	enumDefaultValue = "UNSPECIFIED"
)

const (
	PayloadType_Namespace = "namespace"
	PayloadType_Class     = "class"
	PayloadType_Enum      = "enum"
)

// RE expression
const (
	importREExp       = `import\s*"(\w+)\.proto"\s*;`                   // e.g. 'import "common.proto";', get 'common'
	messageStartREExp = `message\s+(\w+)\s*{`                           // e.g. 'message xxx {', get 'xxx'
	messageFieldREExp = `(repeated)?\s*([\w\.]+)\s+(\w+)\s*=\s*\d+\s*;` // e.g. 'common.Error err = 1;', get 'common.Error'/'err'/'1'
	enumStartREExp    = `enum\s+(\w+)\s*{`                              // e.g. 'enum xxx {', get 'xxx'
	enumFieldREExp    = `(\w+)\s*=\s*\d+\s*;`                           // e.g. 'UNSPECIFIED = 0;', get 'UNSPECIFIED'/'0'
	endREExp          = `}`                                             // end of 'message' and 'enum'
)

type REMatched = uint8

const (
	REMatched_UnMatched REMatched = iota
	REMatched_Import
	REMatched_MessageStart
	REMatched_MessageField
	REMatched_EnumStart
	REMatched_EnumField
	REMatched_End
)
