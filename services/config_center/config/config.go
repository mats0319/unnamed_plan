package config

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
	"os"
)

const uid_ConfigCenterServiceConfig = "8066a746-d419-46df-b213-69af9ea4deac"

type configCenterServiceConfig struct {
}

var configCenterServiceConfigIns = &configCenterServiceConfig{}

func init() {
	byteSlice := mconfig.GetConfig(uid_ConfigCenterServiceConfig)

	err := json.Unmarshal(byteSlice, configCenterServiceConfigIns)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", uid_ConfigCenterServiceConfig), zap.Error(err))
		os.Exit(-1)
	}

	mlog.Logger().Info("> User service config init finish.")
}

func GetConfig() *configCenterServiceConfig {
	return configCenterServiceConfigIns
}
