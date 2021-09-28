package http

import (
	"net/http"

	mconfig "github.com/mats9693/utils/toy_server/config"
	mconst "github.com/mats9693/utils/toy_server/const"
)

var isDev bool

func init() {
	isDev = mconfig.GetConfigLevel() == mconst.ConfigDevLevel

	// user
	http.HandleFunc("/api/login", login)

	// cloud file
	http.HandleFunc("/api/cloudFile/listByUploader", listCloudFileByUploader)
	http.HandleFunc("/api/cloudFile/listPublic", listPublicCloudFile)
}
