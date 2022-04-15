package main

import (
	"github.com/mats9693/unnamed_plan/services/gateway/http"
	"github.com/mats9693/unnamed_plan/services/gateway/rpc"
	"github.com/mats9693/unnamed_plan/services/shared/http"
	"github.com/mats9693/unnamed_plan/services/shared/init"
)

func main() {
	initialize.Init("config.json", http.Init, rpc.Init)

	mhttp.StartServer(http.GetHandler())
}
