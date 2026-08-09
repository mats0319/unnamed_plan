package mhttp

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

type Context struct {
	writer  http.ResponseWriter
	request *http.Request

	AccessToken string // 登录成功后签发，需要身份认证的接口应包含此项
	UserName    string // parse from 'AccessToken'

	ResData any // expect: *utils.Error / *api.resStruct
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		writer:      w,
		request:     r,
		AccessToken: r.Header.Get(utils.HTTPHeader_AccessToken),
	}
}

// ParseParams 这个函数读取了http req.body，这一结构被限制只能读取一次，
// 所以包括请求处理函数和中间件在内，如果有多个函数均调用该函数，会出现错误。
func (ctx *Context) ParseParams(obj any) bool {
	bodyBytes, err := io.ReadAll(ctx.request.Body)
	if err != nil {
		e := utils.ErrServerInternalError().WithCause(err)
		mlog.Error(e.String())
		ctx.ResData = e
		return false
	}

	err = json.Unmarshal(bodyBytes, obj)
	if err != nil {
		e := utils.ErrDeserializeReqParam().WithCause(err)
		mlog.Error(e.String())
		ctx.ResData = e
		return false
	}

	return true
}

func (ctx *Context) SetHeader(key string, value string) {
	ctx.writer.Header().Set(key, value)
}

type Response struct {
	IsSuccess bool   `json:"is_success"`
	Code      int    `json:"code"`
	Err       string `json:"err"`
	Data      any    `json:"data"`
}

// response 该函数不应该中途返回，一定要执行到write
func (ctx *Context) response() {
	httpCode := http.StatusOK

	var obj any
	switch v := ctx.ResData.(type) {
	case *utils.Error:
		obj = &Response{Code: v.Code, Err: v.Error()}

		httpCode = v.HTTPCode
	default: // *api.resStruct(s)
		obj = &Response{IsSuccess: true, Data: v}
	}

	resBytes, err := json.Marshal(obj)
	if err != nil {
		mlog.Error("serialize res to json failed", slog.Any("error", err))
	}

	// write res
	ctx.writer.WriteHeader(httpCode)

	_, err = ctx.writer.Write(resBytes)
	if err != nil {
		mlog.Error("response failed", slog.Any("error", err))
		return
	}
}
