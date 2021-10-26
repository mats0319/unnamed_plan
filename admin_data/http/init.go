package http

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/admin_data/config"
	"github.com/mats9693/unnamed_plan/admin_data/http/handlers"
	"github.com/mats9693/utils/toy_server/http"
	"github.com/mats9693/utils/toy_server/utils"
	"os"
)

var Handlers *mhttp.Handlers

func init() {
	Handlers = mhttp.NewHandlers()

	// user
	Handlers.HandleFunc("/api/login", handlers.Login)
	Handlers.HandleFunc("/api/user/list", handlers.ListUser)
	Handlers.HandleFunc("/api/user/create", handlers.CreateUser)
	Handlers.HandleFunc("/api/user/lock", handlers.LockUser)
	Handlers.HandleFunc("/api/user/unlock", handlers.UnlockUser)
	Handlers.HandleFunc("/api/user/modifyInfo", handlers.ModifyUserInfo)
	Handlers.HandleFunc("/api/user/modifyPermission", handlers.ModifyUserPermission)

	// cloud file
	initCloudFileDir()
	Handlers.HandleFunc("/api/cloudFile/listByUploader", handlers.ListCloudFileByUploader)
	Handlers.HandleFunc("/api/cloudFile/listPublic", handlers.ListPublicCloudFile)
	Handlers.HandleFunc("/api/cloudFile/upload", handlers.UploadCloudFile)
	Handlers.HandleFunc("/api/cloudFile/modify", handlers.ModifyCloudFile)
	Handlers.HandleFunc("/api/cloudFile/delete", handlers.DeleteCloudFile)

	// thinking note
	Handlers.HandleFunc("/api/thinkingNote/listByWriter", handlers.ListThinkingNoteByWriter)
	Handlers.HandleFunc("/api/thinkingNote/listPublic", handlers.ListPublicThinkingNote)
	Handlers.HandleFunc("/api/thinkingNote/create", handlers.CreateThinkingNote)
	Handlers.HandleFunc("/api/thinkingNote/modify", handlers.ModifyThinkingNote)
	Handlers.HandleFunc("/api/thinkingNote/delete", handlers.DeleteThinkingNote)
}

func initCloudFileDir() {
	root := system_config.GetConfiguration().CloudFileRootPath
	path := mutils.FormatDirSuffix(root) + system_config.GetConfiguration().CloudFilePublicDir

	err := os.MkdirAll(path, 0755)
	if err != nil {
		fmt.Println("os.MkdirAll failed, error:", err.Error())
		os.Exit(-1)
	}

	fmt.Println("> Cloud file directory init finish.")
}
