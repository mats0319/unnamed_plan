package http

import (
	"github.com/mats9693/unnamed_plan/shared/go/config"
	. "github.com/mats9693/unnamed_plan/shared/go/const"
	"net/http"
)

var isDev bool

func init() {
	isDev = config.GetConfigLevel() == ConfigDevLevel

	// todo: 在http请求拦截器中，为所有的处理函数添加一些公共处理，例如开发模式允许跨域、http请求参数反序列化等
	http.HandleFunc("/api/login", login)
	http.HandleFunc("/api/user/list", listUser)
	http.HandleFunc("/api/user/create", createUser)
	http.HandleFunc("/api/user/lock", lockUser)
	http.HandleFunc("/api/user/unlock", unlockUser)
	http.HandleFunc("/api/user/modifyPermission", modifyUserPermission)
	http.HandleFunc("/api/user/modifyInfo", modifyUserInfo)
}
