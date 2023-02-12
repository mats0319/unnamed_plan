package mhttp

import (
	"net/http"
)

// Plugins
// stop: if a plugin runs wrong, stop the request or not is according to it
type Plugins interface {
	RunBeforeHook(uri string) bool
	RunAfterHook(uri string) bool

	BeforeHook(request *http.Request, timestamp int64) (code int, err error, stop bool)
	AfterHook(writer http.ResponseWriter, request *http.Request, timestamp int64, param string) (code int, err error, stop bool)
}
