package http

import (
    "github.com/mats9693/unnamed_plan/services/config_center/http/handlers"
    "github.com/mats9693/unnamed_plan/services/shared/const"
    "github.com/mats9693/unnamed_plan/services/shared/http"
)

var handlersIns *mhttp.Handlers

func GetHandler() *mhttp.Handlers {
    return handlersIns
}

func init() {
    handlersIns = mhttp.NewHandlers()

    // authenticate
    handlersIns.HandleFunc("/api/cc/login", handlers.Login, mconst.SkipLimit, mconst.ReSetParams)
}
