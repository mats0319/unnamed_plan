package config

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
	"os"
)

const uid_TaskServiceConfig = "034ac632-8329-4045-94a6-92219b22a263"

type taskServiceConfig struct {
	Address           string `json:"address"`
	UserServerAddress string `json:"userServerAddress"`
	MaxRecords        int    `json:"maxRecords"`
}

var taskServiceConfigIns = &taskServiceConfig{}

func init() {
	byteSlice := mconfig.GetConfig(uid_TaskServiceConfig)

	err := json.Unmarshal(byteSlice, taskServiceConfigIns)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", uid_TaskServiceConfig), zap.Error(err))
		os.Exit(-1)
	}

	mlog.Logger().Info("> Task service config init finish.")
}

func GetConfig() *taskServiceConfig {
	return taskServiceConfigIns
}
