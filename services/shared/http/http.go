package mhttp

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/http/plugins/limit_multi_login"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
	"net/http"
)

type httpConfig struct {
	Port             string   `json:"port"`
	Sources          []string `json:"sources"`
	UnlimitedSources []string `json:"unlimitedSources"`

	LimitMultiLoginConfig limit_multi_login.LimitMultiLoginConfig `json:"limitMultiLoginConfig"`
}

// StartServer is blocked
func StartServer(handlers *Handlers) {
	handlers.config = getHttpConfig()

	handlers.isDev = mconfig.GetConfigLevel() == mconst.ConfigLevel_Dev

	if handlers.config.LimitMultiLoginConfig.LimitMultiLogin {
		handlers.plugins = append(handlers.plugins, limit_multi_login.Init(&handlers.config.LimitMultiLoginConfig,
			handlers.config.Sources))
	}

	mlog.Logger().Info("> Listening at : " + handlers.config.Port)

	// blocked
	err := http.ListenAndServe(":"+handlers.config.Port, handlers)
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
