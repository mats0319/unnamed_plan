package config

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
	"os"
)

const uid_CloudFileServiceConfig = "76175dee-ffee-4198-8f97-94dfd4f960fd"

type cloudFileServiceConfig struct {
	init bool

	Address            string `json:"address"`
	UserServerAddress  string `json:"userServerAddress"`
	CloudFileRootPath  string `json:"cloudFileRootPath"`  // absolute path
	CloudFilePublicDir string `json:"cloudFilePublicDir"` // public folder name
}

var cloudFileServiceConfigIns = &cloudFileServiceConfig{}

func Init() {
	if cloudFileServiceConfigIns.init { // have initialized
		return
	}

	byteSlice := mconfig.GetConfig(uid_CloudFileServiceConfig)

	err := json.Unmarshal(byteSlice, cloudFileServiceConfigIns)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", uid_CloudFileServiceConfig), zap.Error(err))
		os.Exit(-1)
	}

	cloudFileServiceConfigIns.init = true

	mlog.Logger().Info("> Cloud file service config init finish.")
}

func GetConfig() *cloudFileServiceConfig {
	return cloudFileServiceConfigIns
}
