package mconst

type HTTPMultiLoginFlag = uint8

const (
	HTTPMultiLogin_SkipLimit HTTPMultiLoginFlag = iota
	HTTPMultiLogin_ReSetParams
)

type TaskStatus = uint8

const (
	TaskStatus_History = iota // history task will not display
	TaskStatus_Posted
	TaskStatus_InProgress
	TaskStatus_Completed
)
