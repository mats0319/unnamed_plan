package utils

import (
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
)

type Error struct {
	Code    uint32
	Message string
}

const (
	dbErrorCode        uint32 = 10001
	executionErrorCode uint32 = 10004
)

var (
	Error_InvalidParams            = &Error{Code: 10000, Message: mconst.Error_InvalidParams}
	Error_InvalidAccountOrPassword = &Error{Code: 10002, Message: mconst.Error_InvalidAccountOrPassword}
	Error_PermissionDenied         = &Error{Code: 10003, Message: mconst.Error_PermissionDenied}
	Error_NoValidModification      = &Error{Code: 10005, Message: mconst.Error_NoValidModification}
)

func (e *Error) ToRPC() *rpc_impl.Error {
	return &rpc_impl.Error{
		Code:    e.Code,
		Message: e.Message,
	}
}

func NewDBError(message string) *Error {
	return &Error{
		Code:    dbErrorCode,
		Message: message,
	}
}

func NewExecError(message string) *Error {
	return &Error{
		Code:    executionErrorCode,
		Message: message,
	}
}
