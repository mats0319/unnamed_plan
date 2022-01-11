package mhttp

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type handler func(r *http.Request) *ResponseData

type handlerWrapper struct {
	handler handler

	skipMultiLoginLimit   bool
	reSetMultiLoginParams bool
}

type Handlers struct {
	handlersMap sync.Map // pattern - handlerWrapper

	config *httpConfig
	isDev  bool

	loginInfoMap sync.Map // user id - login info
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

	// parse req source
	source := request.Header.Get(mconst.HTTP_MultiLoginSourceSign)

	isLimitedSource := utils.Contains(h.config.Sources, source)
	if !isLimitedSource && utils.Contains(h.config.UnlimitedSources, source) {
		http.Error(writer, mconst.Error_UnknownSource+source, http.StatusUnauthorized)
		return
	}

	// todo: 尝试获取访问来源ip，参考nginx日志
	mlog.Logger().Info("> Receive new request:", zap.String("uri", request.RequestURI))

	// get matched handler by URI
	wrapper := &handlerWrapper{}
	{
		v, ok := h.handlersMap.Load(request.RequestURI)
		if !ok {
			http.Error(writer, mconst.Error_UnknownURI+request.RequestURI, http.StatusNotImplemented)
			return
		}

		wrapper, ok = v.(*handlerWrapper)
		if !ok {
			http.Error(writer, "", http.StatusInternalServerError)
			return
		}
	}

	// multi-login limit: verify
	timestamp := time.Now().Unix()
	if h.config.LimitMultiLogin && isLimitedSource && !wrapper.skipMultiLoginLimit {
		userID := request.Header.Get(mconst.HTTP_MultiLoginUserIDSign)
		token := request.Header.Get(mconst.HTTP_MultiLoginTokenSign)

		if errMsg := h.multiLoginVerify(userID, source, token, timestamp); len(errMsg) > 0 {
			http.Error(writer, errMsg, http.StatusUnauthorized)
			return
		}
	}

	// invoke handler func
	res := wrapper.handler(request)

	// multi-login limit: refresh token
	if h.config.LimitMultiLogin && isLimitedSource && wrapper.reSetMultiLoginParams && !res.HasError {
		newToken := utils.RandomHexString(10)

		res.Token = newToken
		if errMsg := h.setLoginInfo(res.UserID, source, newToken, timestamp); len(errMsg) > 0 {
			http.Error(writer, errMsg, http.StatusInternalServerError)
			return
		}
	}

	h.response(writer, res)

	mlog.Logger().Info("> Handle request result:", zap.Bool("success", res.HasError),
		zap.String("uri", request.RequestURI))

	return
}

// HandleFunc register pattern - handler pair into http handlers
//
// @params params: describe if this handler has different behaviors from default
func (h *Handlers) HandleFunc(pattern string, handler handler, params ...mconst.MultiLoginSign) {
	handlerWrapperIns := &handlerWrapper{
		handler: handler,
	}

	for i := range params {
		switch params[i] {
		case mconst.SkipLimit:
			handlerWrapperIns.skipMultiLoginLimit = true
		case mconst.ReSetParams:
			handlerWrapperIns.reSetMultiLoginParams = true
		}
	}

	h.handlersMap.Store(pattern, handlerWrapperIns)
}

func (h *Handlers) response(writer http.ResponseWriter, data *ResponseData) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(jsonBytes)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}
