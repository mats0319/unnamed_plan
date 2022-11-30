package http

import (
	"github.com/mats9693/unnamed_plan/services/gateway/http/handlers"
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
	handlersIns.HandleFunc("/api/user/login", handlers.Login)
	handlersIns.HandleFunc("/api/user/list", handlers.ListUser)
	handlersIns.HandleFunc("/api/user/create", handlers.CreateUser)
	handlersIns.HandleFunc("/api/user/lock", handlers.LockUser)
	handlersIns.HandleFunc("/api/user/unlock", handlers.UnlockUser)
	handlersIns.HandleFunc("/api/user/modifyInfo", handlers.ModifyUserInfo)
	handlersIns.HandleFunc("/api/user/modifyPermission", handlers.ModifyUserPermission)

	// cloud file
	handlersIns.HandleFunc("/api/cloudFile/list", handlers.ListCloudFile)
	handlersIns.HandleFunc("/api/cloudFile/upload", handlers.UploadCloudFile)
	handlersIns.HandleFunc("/api/cloudFile/modify", handlers.ModifyCloudFile)
	handlersIns.HandleFunc("/api/cloudFile/delete", handlers.DeleteCloudFile)

	// note
	handlersIns.HandleFunc("/api/note/list", handlers.ListNote)
	handlersIns.HandleFunc("/api/note/create", handlers.CreateNote)
	handlersIns.HandleFunc("/api/note/modify", handlers.ModifyNote)
	handlersIns.HandleFunc("/api/note/delete", handlers.DeleteNote)

	// plugins
	limit_multi_login.HandleFunc("/api/login", mconst.HTTPMultiLogin_SkipLimit, mconst.HTTPMultiLogin_ReSetParams)

	return nil
}
