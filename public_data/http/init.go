package http

import "net/http"

func init() {
	http.HandleFunc("/api/login", login)
}
