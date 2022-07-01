package mhttp

import (
	"net/http"
)

type Plugins interface {
	BeforeInvokeHook(request *http.Request, timestamp int64) (code int, err error)
	AfterInvokeHook(response *ResponseData, request *http.Request, timestamp int64) (code int, err error)
}
