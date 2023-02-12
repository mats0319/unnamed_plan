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

	// allow self-define headers in req and res
	writer.Header().Set("Access-Control-Allow-Headers", "*")
	writer.Header().Set("Access-Control-Expose-Headers", "*")

	// verify request
	if stop := h.verifyRequest(writer, request); stop {
		return
	}

	// log req
	timestamp := time.Now().Unix()
	mlog.Logger().Info("> Receive new request",
		zap.String("uri", request.RequestURI),
		zap.String("remote address", request.RemoteAddr),
		zap.Int64("timestamp", timestamp))

	var res *ResponseData

	// log res
	defer func() {
		fields := make([]zap.Field, 0, 3)
		fields = append(fields, zap.String("uri", request.RequestURI))
		fields = append(fields, zap.Bool("is success", res.Success))

		if !res.Success {
			fields = append(fields, zap.String("error", string(res.Data)))
		}

		mlog.Logger().Info("> Handle request", fields...)
	}()

	// plugins: before invoke hooks
	if stop := h.runBeforeInvokeHooks(writer, request, timestamp); stop {
		return
	}

	// invoke handleFunc func
	{

		v, ok := h.handlersMap.Load(request.RequestURI)
		if !ok {
			http.Error(writer, mconst.Error_UnknownURI+request.RequestURI, http.StatusNotImplemented)
			return
		}

		handleFuncIns, ok := v.(handleFunc)
		if !ok {
			http.Error(writer, "type assert error", http.StatusInternalServerError)
			return
		}

		res = handleFuncIns(request)
	}

	// plugins: after invoke hooks
	if stop := h.runAfterInvokeHooks(writer, request, res, timestamp); stop {
		return
	}

	// response
	response(writer, res.Data)
}

// HandleFunc register pattern - handleFunc pair into http handlers
func (h *Handlers) HandleFunc(pattern string, hf handleFunc) {
	h.handlersMap.Store(pattern, hf)
}

func (h *Handlers) verifyRequest(writer http.ResponseWriter, request *http.Request) bool {
	if request.Method != http.MethodPost {
		response(writer, []byte(""))
		return true
	}

	source := request.Header.Get(mconst.HTTP_SourceSign)
	if !utils.Contains(h.config.Sources, source) {
		http.Error(writer, mconst.Error_UnknownSource+source, http.StatusUnauthorized)
		return true
	}

	return false
}

func (h *Handlers) runBeforeInvokeHooks(writer http.ResponseWriter, request *http.Request, timestamp int64) bool {
	var (
		code int
		err  error
		stop bool
	)
	for i := range h.plugins {
		if runHook := h.plugins[i].RunBeforeHook(request.RequestURI); !runHook {
			continue
		}

		code, err, stop = h.plugins[i].BeforeHook(request, timestamp)
		if err != nil {
			mlog.Logger().Error("> Plugin: before invoke hook exec error",
				zap.Int("plugin index", i),
				zap.Int("error code", code),
				zap.Error(err))
			if stop {
				http.Error(writer, err.Error(), code)
				break
			}
		}
	}

	return stop
}

func (h *Handlers) runAfterInvokeHooks(writer http.ResponseWriter, request *http.Request, res *ResponseData, timestamp int64) bool {
	if !res.Success {
		return false
	}

	param := ""
	if res != nil {
		param = res.AfterHookParam
	}

	var (
		code int
		err  error
		stop bool
	)
	for i := range h.plugins {
		if runHook := h.plugins[i].RunAfterHook(request.RequestURI); !runHook {
			continue
		}

		code, err, stop = h.plugins[i].AfterHook(writer, request, timestamp, param)
		if err != nil {
			mlog.Logger().Error("> Plugin: after invoke hook exec error",
				zap.Int("plugin index", i),
				zap.Int("error code", code),
				zap.Error(err))
			if stop {
				http.Error(writer, err.Error(), code)
				break
			}
		}
	}

	return stop
}
