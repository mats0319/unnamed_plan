package mconst

type HTTPFlags = uint8

const (
	HTTPMultiLogin_SkipLimit HTTPFlags = iota
	HTTPMultiLogin_ReSetParams
)

type TaskStatus = uint8

const (
	TaskStatus_History = iota // history task will not display
	TaskStatus_Posted
	TaskStatus_InProgress
	TaskStatus_Completed
)

type ConfigLevel string

const (
	ConfigLevel_Default = "default"
	ConfigLevel_Dev = "dev"
	ConfigLevel_Product = "product"
	ConfigLevel_Test = "test"
)
