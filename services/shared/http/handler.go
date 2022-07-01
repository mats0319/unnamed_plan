package mhttp

import (
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type handleFunc func(r *http.Request) *ResponseData

type Handlers struct {
	isDev  bool
	config *httpConfig

	handlersMap sync.Map // pattern - handleFunc

	plugins []Plugins
}

var _ http.Handler = (*Handlers)(nil)

func NewHandlers() *Handlers {
	return &Handlers{
		handlersMap: sync.Map{},
		plugins:     make([]Plugins, 0),
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
		if request.Method != http.MethodPost {
			response(writer, &ResponseData{})
			return
		}

		source := request.Header.Get(mconst.HTTP_SourceSign)
		if !utils.Contains(h.config.Sources, source) && !utils.Contains(h.config.UnlimitedSources, source) {
			http.Error(writer, mconst.Error_UnknownSource+source, http.StatusUnauthorized)
			return
		}
	}

	// log
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

	// invoke handleFunc func
	var res *ResponseData
	{
		v, ok := h.handlersMap.Load(request.RequestURI)
		if !ok {
			http.Error(writer, mconst.Error_UnknownURI+request.RequestURI, http.StatusNotImplemented)
			return
		}

		handlerFunc, ok := v.(handleFunc)
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

	// return
	response(writer, res)

	// log
	logData := []zap.Field{
		zap.String("uri", request.RequestURI),
		zap.Bool("is success", !res.HasError),
	}

	if res.HasError {
		logData = append(logData, zap.String("error", res.Data))
	}

	mlog.Logger().Info("> Handle request result:", logData...)

	return
}

// HandleFunc register pattern - handleFunc pair into http handlers
func (h *Handlers) HandleFunc(pattern string, hf handleFunc) {
	h.handlersMap.Store(pattern, hf)
}
