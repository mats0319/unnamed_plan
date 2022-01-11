package mhttp

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type httpConfig struct {
	Port             string   `json:"port"`
	LimitMultiLogin  bool     `json:"limitMultiLogin"`
	KeepTokenValid   int64    `json:"KeepTokenValid"` // unit: second
	Sources          []string `json:"sources"`
	UnlimitedSources []string `json:"unlimitedSources"`
}

// StartServer is blocked
func StartServer(handlers *Handlers) {
	handlers.config = getHttpConfig()

	handlers.isDev = mconfig.GetConfigLevel() == mconst.ConfigDevLevel

	if handlers.config.LimitMultiLogin {
		handlers.loginInfoMap = sync.Map{}
	}

	mlog.Logger().Info("> Listening at : " + handlers.config.Port)

	// blocked
	err := http.ListenAndServe(":"+handlers.config.Port, handlers)
	if err != nil {
		mlog.Logger().Error("http listen and serve failed", zap.Error(err))
	}
}

func getHttpConfig() *httpConfig {
	byteSlice := mconfig.GetConfig(mconst.UID_HTTP)

	conf := &httpConfig{}
	err := json.Unmarshal(byteSlice, conf)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", mconst.UID_HTTP), zap.Error(err))
		return nil
	}

	return conf
}
