package http

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/admin_data/config"
	"github.com/mats9693/unnamed_plan/admin_data/http/handlers"
	"github.com/mats9693/unnamed_plan/admin_data/kits"
	"net/http"
	"os"
)

func init() {
	// todo: 在http请求拦截器中，为所有的处理函数添加一些公共处理，例如开发模式允许跨域、http请求参数反序列化等
	// user
	http.HandleFunc("/api/login", handlers.Login)
	http.HandleFunc("/api/user/list", handlers.ListUser)
	http.HandleFunc("/api/user/create", handlers.CreateUser)
	http.HandleFunc("/api/user/lock", handlers.LockUser)
	http.HandleFunc("/api/user/unlock", handlers.UnlockUser)
	http.HandleFunc("/api/user/modifyInfo", handlers.ModifyUserInfo)
	http.HandleFunc("/api/user/modifyPermission", handlers.ModifyUserPermission)

	// cloud file
	initCloudFileDir()
	http.HandleFunc("/api/cloudFile/listByUploader", handlers.ListCloudFileByUploader)
	http.HandleFunc("/api/cloudFile/listPublic", handlers.ListPublicCloudFile)
	http.HandleFunc("/api/cloudFile/upload", handlers.UploadCloudFile)
	http.HandleFunc("/api/cloudFile/modify", handlers.ModifyCloudFile)
	http.HandleFunc("/api/cloudFile/delete", handlers.DeleteCloudFile)

	// thinking note
	http.HandleFunc("/api/thinkingNote/listByWriter", handlers.ListThinkingNoteByWriter)
	http.HandleFunc("/api/thinkingNote/listPublic", handlers.ListPublicThinkingNote)
	http.HandleFunc("/api/thinkingNote/create", handlers.CreateThinkingNote)
	http.HandleFunc("/api/thinkingNote/modify", handlers.ModifyThinkingNote)
	http.HandleFunc("/api/thinkingNote/delete", handlers.DeleteThinkingNote)
}

func initCloudFileDir() {
	root := system_config.GetConfiguration().CloudFileRootPath
	path := kits.AppendDirSuffix(root) + system_config.GetConfiguration().CloudFilePublicDir

	err := os.MkdirAll(path, 0755)
	if err != nil {
		fmt.Println("os.MkdirAll failed, error:", err.Error())
		os.Exit(-1)
	}

	fmt.Println("> Cloud file directory init finish.")
}
