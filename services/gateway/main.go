package main

import (
	_ "github.com/mats9693/utils/toy_server/config"

	"github.com/mats9693/unnamed_plan/services/gateway/http"
	"github.com/mats9693/utils/toy_server/http"
)

func main() {
	mhttp.StartServer(http.Handlers)
}