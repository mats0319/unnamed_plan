package http

import (
	"github.com/mats9693/unnamed_plan/services/gateway/plugins/limit_multi_login"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/http"
)

var handlersIns *mhttp.Handlers

func GetHandler() *mhttp.Handlers {
	return handlersIns
}

func Init() error {
	handlersIns = mhttp.NewHandlers()

	// user
	handlersIns.HandleFunc("/api/user/login", Login)
	handlersIns.HandleFunc("/api/user/list", ListUser)
	handlersIns.HandleFunc("/api/user/create", CreateUser)
	handlersIns.HandleFunc("/api/user/lock", LockUser)
	handlersIns.HandleFunc("/api/user/unlock", UnlockUser)
	handlersIns.HandleFunc("/api/user/modifyInfo", ModifyUserInfo)
	handlersIns.HandleFunc("/api/user/modifyPermission", ModifyUserPermission)

	//cloud file
	handlersIns.HandleFunc("/api/cloudFile/list", ListCloudFile)
	handlersIns.HandleFunc("/api/cloudFile/upload", UploadCloudFile)
	handlersIns.HandleFunc("/api/cloudFile/modify", ModifyCloudFile)
	handlersIns.HandleFunc("/api/cloudFile/delete", DeleteCloudFile)

	// plugins
	limit_multi_login.HandleFunc("/api/user/login", mconst.HTTPFlags_MultiLogin_SkipLimit, mconst.HTTPFlags_MultiLogin_ReSetParams)

	return nil
}
