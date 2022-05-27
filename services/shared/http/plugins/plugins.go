package plugins

import (
	"github.com/mats9693/unnamed_plan/services/shared/http/response"
	"net/http"
)

type Plugins interface {
	BeforeInvokeHook(request *http.Request, timestamp int64) (code int, err error)
	AfterInvokeHook(response *mresponse.ResponseData, request *http.Request, timestamp int64) (code int, err error)
}
