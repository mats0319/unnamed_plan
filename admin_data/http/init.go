package http

import (
	"github.com/mats9693/unnamed_plan/shared/go/config"
	. "github.com/mats9693/unnamed_plan/shared/go/const"
	"net/http"
)

var isDev bool

func init() {
	isDev = config.GetConfigLevel() == ConfigDevLevel

	http.HandleFunc("/api/login", login)
	http.HandleFunc("/api/user/list", listUser)
	http.HandleFunc("/api/user/create", createUser)
}
