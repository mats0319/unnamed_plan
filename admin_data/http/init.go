package http

import (
	"net/http"

	mconfig "github.com/mats9693/utils/toy_server/config"
	mconst "github.com/mats9693/utils/toy_server/const"
)

var isDev bool

func init() {
	isDev = mconfig.GetConfigLevel() == mconst.ConfigDevLevel

	// todo: 在http请求拦截器中，为所有的处理函数添加一些公共处理，例如开发模式允许跨域、http请求参数反序列化等
	// user
	http.HandleFunc("/api/login", login)
	http.HandleFunc("/api/user/list", listUser)
	http.HandleFunc("/api/user/create", createUser)
	http.HandleFunc("/api/user/lock", lockUser)
	http.HandleFunc("/api/user/unlock", unlockUser)
	http.HandleFunc("/api/user/modifyPermission", modifyUserPermission)
	http.HandleFunc("/api/user/modifyInfo", modifyUserInfo)

	// cloud file
	http.HandleFunc("/api/cloudFile/upload", uploadFile)
	http.HandleFunc("/api/cloudFile/listByUploader", listCloudFileByUploader)
	http.HandleFunc("/api/cloudFile/listPublic", listPublicCloudFile)
	http.HandleFunc("/api/cloudFile/listDelete", listDeleted)
	http.HandleFunc("/api/cloudFile/delete", deleteFile)
}
