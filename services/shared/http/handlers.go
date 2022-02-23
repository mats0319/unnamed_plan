package mhttp

import (
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/http/plugins"
	"github.com/mats9693/unnamed_plan/services/shared/http/plugins/limit_multi_login"
	"github.com/mats9693/unnamed_plan/services/shared/http/response"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type handler func(r *http.Request) *mresponse.ResponseData

type Handlers struct {
	isDev  bool
	config *httpConfig

	handlersMap sync.Map // pattern - handler

	plugins []plugins.Plugins
}

var _ http.Handler = (*Handlers)(nil)

func NewHandlers() *Handlers {
	return &Handlers{
		handlersMap: sync.Map{},
	}
}

func (h *Handlers) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// allow dev mode cross-origin
	if h.isDev {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
	}

	// allow self-define headers
	writer.Header().Set("Access-Control-Allow-Headers", "*")

	// verify request
	{
		if request.Method == http.MethodOptions {
			response(writer, &mresponse.ResponseData{})
			return
		}

		source := request.Header.Get(mconst.HTTP_SourceSign)
		if !utils.Contains(h.config.Sources, source) && !utils.Contains(h.config.UnlimitedSources, source) {
			http.Error(writer, mconst.Error_UnknownSource+source, http.StatusUnauthorized)
			return
		}
	}

	timestamp := time.Now().Unix()
	mlog.Logger().Info("> Receive new request:",
		zap.String("uri", request.RequestURI),
		zap.String("remote address", request.RemoteAddr),
		zap.Int64("timestamp", timestamp))

	// plugins: before invoke hooks
	for i := range h.plugins {
		code, err := h.plugins[i].BeforeInvokeHook(request, timestamp)
		if err != nil { // make sure hook logged each error before return
			http.Error(writer, err.Error(), code)
			return
		}
	}

	// invoke handler func
	var res *mresponse.ResponseData
	{
		v, ok := h.handlersMap.Load(request.RequestURI)
		if !ok {
			http.Error(writer, mconst.Error_UnknownURI+request.RequestURI, http.StatusNotImplemented)
			return
		}

		handlerFunc, ok := v.(handler)
		if !ok {
			http.Error(writer, "type assert error", http.StatusInternalServerError)
			return
		}

		res = handlerFunc(request)
	}

	// plugins: after invoke hooks
	if !res.HasError {
		for i := range h.plugins {
			code, err := h.plugins[i].AfterInvokeHook(res, request, timestamp)
			if err != nil { // make sure hook logged each error before return
				http.Error(writer, err.Error(), code)
				return
			}
		}
	}

	response(writer, res)

	mlog.Logger().Info("> Handle request result:",
		zap.String("uri", request.RequestURI),
		zap.Bool("success", !res.HasError))

	return
}

// HandleFunc register pattern - handler pair into http handlers
//
//   @params params: http plugins flags
func (h *Handlers) HandleFunc(pattern string, handlerFunc handler, params ...uint8) {
	h.handlersMap.Store(pattern, handlerFunc)

	limitMultiLoginParams := make([]uint8, 0, 2)
	for i := range params {
		switch params[i] {
		case mconst.HTTPMultiLogin_SkipLimit, mconst.HTTPMultiLogin_ReSetParams:
			limitMultiLoginParams = append(limitMultiLoginParams, params[i])
		}
	}

	// even params is empty, just register pattern, distinguish from absolutely unknown uri
	limit_multi_login.HandleFunc(pattern, limitMultiLoginParams...)
}
