package main

import (
	"github.com/mats9693/unnamed_plan/services/gateway/http"
	"github.com/mats9693/unnamed_plan/services/gateway/rpc"
	mconst "github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/http"
	"github.com/mats9693/unnamed_plan/services/shared/init"
)

func main() {
	initialize.InitFromConfigCenter(mconst.UID_Service_Gateway, http.Init, rpc.Init)

	mhttp.StartServer(http.GetHandler())
}
