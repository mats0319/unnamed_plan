package mhttp

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
	"net/http"
)

type httpConfig struct {
	Port             string   `json:"port"`
	Sources          []string `json:"sources"`
}

// StartServer is blocked
func StartServer(handlers *Handlers, plugins ...Plugins) {
	handlers.config = getHttpConfig()

	handlers.isDev = mconfig.GetConfigLevel() == mconst.ConfigLevel_Dev ||
		mconfig.GetConfigLevel() == mconst.ConfigLevel_Default

	handlers.plugins = plugins

	mlog.Logger().Info("> Listening at : 127.0.0.1:" + handlers.config.Port)

	// blocked
	err := http.ListenAndServe("127.0.0.1:"+handlers.config.Port, handlers)
	if err != nil {
		mlog.Logger().Error("http listen and serve failed", zap.Error(err))
	}
}

func getHttpConfig() *httpConfig {
	byteSlice := mconfig.GetConfig(mconst.UID_Gateway_HTTP)

	conf := &httpConfig{}
	err := json.Unmarshal(byteSlice, conf)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", mconst.UID_Gateway_HTTP), zap.Error(err))
		return nil
	}

	return conf
}
