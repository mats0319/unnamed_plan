package http

import (
	"github.com/mats9693/utils/toy_server/config"
	. "github.com/mats9693/utils/toy_server/const"
	"net/http"
)

var isDev bool

func init() {
	isDev = config.GetConfigLevel() == ConfigDevLevel

	// user
	http.HandleFunc("/api/login", login)

	// cloud file
	http.HandleFunc("/api/cloudFile/listByUploader", listCloudFileByUploader)
	http.HandleFunc("/api/cloudFile/listPublic", listPublicCloudFile)
}
