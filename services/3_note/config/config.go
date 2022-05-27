package config

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
	"os"
)

type noteServiceConfig struct {
	init bool

	Address           string `json:"address"`
	UserServerAddress string `json:"userServerAddress"`
}

var noteServiceConfigIns = &noteServiceConfig{}

func Init() {
	if noteServiceConfigIns.init { // have initialized
		return
	}

	byteSlice := mconfig.GetConfig(mconst.UID_Service_Note)

	err := json.Unmarshal(byteSlice, noteServiceConfigIns)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", mconst.UID_Service_Note), zap.Error(err))
		os.Exit(-1)
	}

	noteServiceConfigIns.init = true

	mlog.Logger().Info("> Note service config init finish.")
}

func GetConfig() *noteServiceConfig {
	return noteServiceConfigIns
}
