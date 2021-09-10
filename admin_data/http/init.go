package http

import (
	"github.com/mats9693/unnamed_plan/admin_data/config"
	. "github.com/mats9693/unnamed_plan/admin_data/const"
	"net/http"
)

var isDev bool

func init() {
	isDev = config.GetConfigLevel() == ConfigDevLevel

	http.HandleFunc("/api/login", Login)
	http.HandleFunc("/api/user/create", createUser)
}
