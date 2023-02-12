package utils

import (
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/proto/go"
)

// common error code, need more information to explain, 4-digits start with '1'
const (
	dbErrorCode             uint32 = 1000
	executionErrorCode             = 1001
	grpcConnectionErrorCode        = 1002
	getClientErrorCode             = 1010
)

// customized error, 4-digits start with '2'
var (
	Error_InvalidParams            = &rpc_impl.Error{Code: 2000, Message: mconst.Error_InvalidParams}
	Error_InvalidAccountOrPassword = &rpc_impl.Error{Code: 2001, Message: mconst.Error_InvalidAccountOrPassword}
	Error_PermissionDenied         = &rpc_impl.Error{Code: 2002, Message: mconst.Error_PermissionDenied}
	Error_NoValidModification      = &rpc_impl.Error{Code: 2003, Message: mconst.Error_NoValidModification}
)

func NewDBError(message string) *rpc_impl.Error {
	return newError(dbErrorCode, message)
}

func NewExecError(message string) *rpc_impl.Error {
	return newError(executionErrorCode, message)
}

func NewGrpcConnectionError(message string) *rpc_impl.Error {
	return newError(grpcConnectionErrorCode, message)
}

func NewGetClientError(message string) *rpc_impl.Error {
	return newError(getClientErrorCode, message)
}

func newError(code uint32, message string) *rpc_impl.Error {
	return &rpc_impl.Error{
		Code:    code,
		Message: message,
	}
}
