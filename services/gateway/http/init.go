package http

import (
	"github.com/mats9693/unnamed_plan/services/gateway/http/handlers"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/http"
)

var handlersIns *mhttp.Handlers

func GetHandler() *mhttp.Handlers {
	return handlersIns
}

func init() {
	handlersIns = mhttp.NewHandlers()

	// user
	handlersIns.HandleFunc("/api/login", handlers.Login, mconst.HTTPMultiLogin_SkipLimit, mconst.HTTPMultiLogin_ReSetParams)
	handlersIns.HandleFunc("/api/user/list", handlers.ListUser)
	handlersIns.HandleFunc("/api/user/create", handlers.CreateUser)
	handlersIns.HandleFunc("/api/user/lock", handlers.LockUser)
	handlersIns.HandleFunc("/api/user/unlock", handlers.UnlockUser)
	handlersIns.HandleFunc("/api/user/modifyInfo", handlers.ModifyUserInfo)
	handlersIns.HandleFunc("/api/user/modifyPermission", handlers.ModifyUserPermission)

	// cloud file
	handlersIns.HandleFunc("/api/cloudFile/listByUploader", handlers.ListCloudFileByUploader)
	handlersIns.HandleFunc("/api/cloudFile/listPublic", handlers.ListPublicCloudFile)
	handlersIns.HandleFunc("/api/cloudFile/upload", handlers.UploadCloudFile)
	handlersIns.HandleFunc("/api/cloudFile/modify", handlers.ModifyCloudFile)
	handlersIns.HandleFunc("/api/cloudFile/delete", handlers.DeleteCloudFile)

	// note
	handlersIns.HandleFunc("/api/note/listByWriter", handlers.ListNoteByWriter)
	handlersIns.HandleFunc("/api/note/listPublic", handlers.ListPublicNote)
	handlersIns.HandleFunc("/api/note/create", handlers.CreateNote)
	handlersIns.HandleFunc("/api/note/modify", handlers.ModifyNote)
	handlersIns.HandleFunc("/api/note/delete", handlers.DeleteNote)

	// task
	handlersIns.HandleFunc("/api/task/list", handlers.ListTask)
	handlersIns.HandleFunc("/api/task/create", handlers.CreateTask)
	handlersIns.HandleFunc("/api/task/modify", handlers.ModifyTask)
}
