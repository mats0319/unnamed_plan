package config

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
)

type taskServiceConfig struct {
	init bool

	MaxRecords int `json:"maxRecords"`
}

var taskServiceConfigIns = &taskServiceConfig{}

func Init() error {
	if taskServiceConfigIns.init { // have initialized
		mlog.Logger().Error("already initialized")
		return nil
	}

	byteSlice := mconfig.GetConfig(mconst.UID_Service_Task)

	err := json.Unmarshal(byteSlice, taskServiceConfigIns)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", mconst.UID_Service_Task), zap.Error(err))
		return err
	}

	taskServiceConfigIns.init = true

	mlog.Logger().Info("> Task service config init finish.")

	return nil
}

func GetConfig() *taskServiceConfig {
	return taskServiceConfigIns
}
