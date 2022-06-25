package config

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
)

type cloudFileServiceConfig struct {
	init bool

	CloudFileRootPath  string `json:"cloudFileRootPath"`  // absolute path
	CloudFilePublicDir string `json:"cloudFilePublicDir"` // public folder name
}

var cloudFileServiceConfigIns = &cloudFileServiceConfig{}

func Init() error {
	if cloudFileServiceConfigIns.init { // have initialized
		mlog.Logger().Error("already initialized")
		return nil
	}

	byteSlice := mconfig.GetConfig(mconst.UID_Service_Cloud_File)

	err := json.Unmarshal(byteSlice, cloudFileServiceConfigIns)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", mconst.UID_Service_Cloud_File), zap.Error(err))
		return utils.NewError(err.Error())
	}

	cloudFileServiceConfigIns.init = true

	mlog.Logger().Info("> Cloud file service config init finish.")

	return nil
}

func GetConfig() *cloudFileServiceConfig {
	return cloudFileServiceConfigIns
}
