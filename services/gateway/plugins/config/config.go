package config

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/gateway/plugins/limit_multi_login"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
)

type gatewayPluginsConfig struct {
	init bool

	LimitMultiLoginConfig *limit_multi_login.LimitMultiLoginConfig `json:"limitMultiLoginConfig"`
}

var gatewayPluginsConfigIns = &gatewayPluginsConfig{}

func Init() error {
	if gatewayPluginsConfigIns.init { // have initialized
		mlog.Logger().Error("already initialized")
		return nil
	}

	byteSlice := mconfig.GetConfig(mconst.UID_Gateway_Plugins)

	err := json.Unmarshal(byteSlice, gatewayPluginsConfigIns)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", mconst.UID_Gateway_Plugins), zap.Error(err))
		return utils.NewError(err.Error())
	}

	// init plugins with config
	{
		if gatewayPluginsConfigIns.LimitMultiLoginConfig != nil {
			limit_multi_login.Init(gatewayPluginsConfigIns.LimitMultiLoginConfig)
		}
	}

	gatewayPluginsConfigIns.init = true

	mlog.Logger().Info("> Gateway plugins config init finish.")

	return nil
}

func GetConfig() *gatewayPluginsConfig {
	return gatewayPluginsConfigIns
}
