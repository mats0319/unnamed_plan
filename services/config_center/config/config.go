package config

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
	"os"
)

type configCenterServiceConfig struct {
	init bool

	Address string `json:"address"`
}

var configCenterServiceConfigIns = &configCenterServiceConfig{}

func Init() {
	if configCenterServiceConfigIns.init { // have initialized
		return
	}

	byteSlice := mconfig.GetConfig(mconst.UID_Service_Config_Center)

	err := json.Unmarshal(byteSlice, configCenterServiceConfigIns)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", mconst.UID_Service_Config_Center), zap.Error(err))
		os.Exit(-1)
	}

	configCenterServiceConfigIns.init = true

	mlog.Logger().Info("> Config Center service config init finish.")
}

func GetConfig() *configCenterServiceConfig {
	return configCenterServiceConfigIns
}
