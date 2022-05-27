package config

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
	"os"
)

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

	byteSlice := mconfig.GetConfig(mconst.UID_Service_Cloud_File)

	err := json.Unmarshal(byteSlice, cloudFileServiceConfigIns)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", mconst.UID_Service_Cloud_File), zap.Error(err))
		os.Exit(-1)
	}

	cloudFileServiceConfigIns.init = true

	mlog.Logger().Info("> Cloud file service config init finish.")
}

func GetConfig() *cloudFileServiceConfig {
	return cloudFileServiceConfigIns
}
