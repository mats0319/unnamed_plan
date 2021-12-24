package http

import (
	"github.com/mats9693/unnamed_plan/services/gateway/http/handlers"
	"github.com/mats9693/utils/toy_server/const"
	"github.com/mats9693/utils/toy_server/http"
)

var handlersIns *mhttp.Handlers

func GetHandler() *mhttp.Handlers {
	return handlersIns
}

func init() {
	handlersIns = mhttp.NewHandlers()

	// user
	handlersIns.HandleFunc("/api/login", handlers.Login, mconst.SkipLimit, mconst.ReSetParams)
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

	// thinking note
	handlersIns.HandleFunc("/api/thinkingNote/listByWriter", handlers.ListThinkingNoteByWriter)
	handlersIns.HandleFunc("/api/thinkingNote/listPublic", handlers.ListPublicThinkingNote)
	handlersIns.HandleFunc("/api/thinkingNote/create", handlers.CreateThinkingNote)
	handlersIns.HandleFunc("/api/thinkingNote/modify", handlers.ModifyThinkingNote)
	handlersIns.HandleFunc("/api/thinkingNote/delete", handlers.DeleteThinkingNote)
}
