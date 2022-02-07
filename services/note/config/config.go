package config

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
	"os"
)

const uid_NoteServiceConfig = "5d71f9af-63f5-44c5-9ef0-403bcd1ad381"

type noteServiceConfig struct {
	Address           string `json:"address"`
	UserServerAddress string `json:"userServerAddress"`
}

var noteServiceConfigIns = &noteServiceConfig{}

func init() {
	byteSlice := mconfig.GetConfig(uid_NoteServiceConfig)

	err := json.Unmarshal(byteSlice, noteServiceConfigIns)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", uid_NoteServiceConfig), zap.Error(err))
		os.Exit(-1)
	}

	mlog.Logger().Info("> Note service config init finish.")
}

func GetConfig() *noteServiceConfig {
	return noteServiceConfigIns
}
