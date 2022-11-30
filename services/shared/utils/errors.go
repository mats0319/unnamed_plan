package utils

import (
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/proto/go"
)

type Error struct {
	Code    uint32
	Message string
}

// common error code, need more information to explain, 4-digits start with '1'
const (
	dbErrorCode             uint32 = 1001
	executionErrorCode      uint32 = 1002
	grpcConnectionErrorCode uint32 = 1003
)

// customized error, 5-digits start with '1'
var (
	Error_InvalidParams            = &Error{Code: 10001, Message: mconst.Error_InvalidParams}
	Error_InvalidAccountOrPassword = &Error{Code: 10002, Message: mconst.Error_InvalidAccountOrPassword}
	Error_PermissionDenied         = &Error{Code: 10003, Message: mconst.Error_PermissionDenied}
	Error_NoValidModification      = &Error{Code: 10004, Message: mconst.Error_NoValidModification}
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

func NewGrpcConnectionError(message string) *Error {
	return &Error{
		Code:    grpcConnectionErrorCode,
		Message: message,
	}
}
