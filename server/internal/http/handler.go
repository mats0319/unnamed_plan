package mhttp

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	mconfig "github.com/mats0319/unnamed_plan/server/internal/config"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

type Handler struct {
	Handlers map[string]*HandlerItem // uri - handler func and middleware(s)
	uri      []string                // 按照添加顺序保存uri，用于启动时展示
}

type HandlerItem struct {
	Func        func(ctx *Context)
	Middlewares []func(ctx *Context) *utils.Error
}

// StartServer is blocked
func (h *Handler) StartServer() {
	h.displayRegisteredURI()

	// 因为允许手机访问需要设置web ip为`192.168.xxx`(即本机内网ipv4地址)，
	// 这里为了支持`127.0.0.1`、内网ip等多种格式，使用`0.0.0.0`
	addr := fmt.Sprintf("0.0.0.0:%d", mconfig.GetInternalConfig().HTTPServerPort)
	mlog.Info("> Listening at: " + addr)

	err := http.ListenAndServe(addr, h)
	if err != nil {
		mlog.Error("handlers listen and serve failed", slog.Any("error", err))
	}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "*")
	writer.Header().Set("Access-Control-Expose-Headers", "*")

	if request.Method != http.MethodPost { // return after set header
		return // res body is empty
	}

	handlerItemIns, ok := h.Handlers[request.RequestURI]
	if !ok {
		// 实际场景中，会触发这一段代码的，更多的是恶意向服务器批量发送请求的，所以就直接返回了。
		// 想看访问细节可以去nginx访问日志。
		// 不想在nginx上配置接口白名单，那样太死板了。找找有什么好办法。
		return
	}

	// log req
	{
		mlog.Info(fmt.Sprintf("> Receive Request: %s .", request.URL.String()))
		startTime := time.Now()
		defer func() {
			mlog.Info(fmt.Sprintf("> Process Request: %s , in %d ms", request.URL.String(), time.Since(startTime).Milliseconds()))

			err := recover()
			if err != nil {
				mlog.Error("recover panic", slog.Any("", err))
			}
		}()
	}

	ctx := NewContext(writer, request)
	defer ctx.response()

	// middlewares
	for i := range handlerItemIns.Middlewares {
		err := handlerItemIns.Middlewares[i](ctx)
		if err != nil { // log in middleware
			ctx.ResData = err
			return
		}
	}

	// business handler
	handlerItemIns.Func(ctx)
}

func (h *Handler) AddHandler(uri string, handlerFunc func(ctx *Context), middlewares ...func(ctx *Context) *utils.Error) {
	if h.Handlers == nil {
		h.Handlers = make(map[string]*HandlerItem)
	}

	h.Handlers[uri] = &HandlerItem{
		Func:        handlerFunc,
		Middlewares: middlewares,
	}

	h.uri = append(h.uri, uri)
}

func (h *Handler) displayRegisteredURI() {
	mlog.Info(fmt.Sprintf("> HTTP Handler Count: %d.", len(h.Handlers)))

	for _, v := range h.uri {
		mlog.Info("- HTTP Handler Registered: " + v)
	}
}
