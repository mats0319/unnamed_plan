package config

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
)

type configCenterServiceConfig struct {
	init bool

	Address string `json:"address"`
}

var configCenterServiceConfigIns = &configCenterServiceConfig{}

func Init() error {
	if configCenterServiceConfigIns.init { // have initialized
		mlog.Logger().Error("already initialized")
		return nil
	}

	byteSlice := mconfig.GetConfig(mconst.UID_Service_Core)

	err := json.Unmarshal(byteSlice, configCenterServiceConfigIns)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", mconst.UID_Service_Core), zap.Error(err))
		return err
	}

	configCenterServiceConfigIns.init = true

	mlog.Logger().Info("> Core service config init finish.")

	return nil
}

func GetConfig() *configCenterServiceConfig {
	return configCenterServiceConfigIns
}
