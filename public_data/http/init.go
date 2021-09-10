package http

import (
	ah "github.com/mats9693/unnamed_plan/admin_data/http"
	"net/http"
)

func init() {
	http.HandleFunc("/api/login", ah.Login)
}
