package http

import (
	"github.com/mats9693/utils/toy_server/config"
	. "github.com/mats9693/utils/toy_server/const"
	"net/http"
)

var isDev bool

func init() {
	isDev = config.GetConfigLevel() == ConfigDevLevel

	http.HandleFunc("/api/login", login)
}
