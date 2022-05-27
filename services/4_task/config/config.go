package config

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
	"os"
)

type taskServiceConfig struct {
	init bool

	Address           string `json:"address"`
	UserServerAddress string `json:"userServerAddress"`
	MaxRecords        int    `json:"maxRecords"`
}

var taskServiceConfigIns = &taskServiceConfig{}

func Init() {
	if taskServiceConfigIns.init { // have initialized
		return
	}

	byteSlice := mconfig.GetConfig(mconst.UID_Service_Task)

	err := json.Unmarshal(byteSlice, taskServiceConfigIns)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", mconst.UID_Service_Task), zap.Error(err))
		os.Exit(-1)
	}

	taskServiceConfigIns.init = true

	mlog.Logger().Info("> Task service config init finish.")
}

func GetConfig() *taskServiceConfig {
	return taskServiceConfigIns
}
