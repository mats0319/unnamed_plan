package config

import (
	"encoding/json"
	"github.com/mats9693/utils/toy_server/config"
	mlog "github.com/mats9693/utils/toy_server/log"
	"go.uber.org/zap"
	"os"
)

const uid_ThinkingNoteServiceConfig = "5d71f9af-63f5-44c5-9ef0-403bcd1ad381"

type thinkingNoteServiceConfig struct {
	Address           string `json:"address"`
	UserServerAddress string `json:"userServerAddress"`
}

var thinkingNoteServiceConfigIns = &thinkingNoteServiceConfig{}

func init() {
	byteSlice := mconfig.GetConfig(uid_ThinkingNoteServiceConfig)

	err := json.Unmarshal(byteSlice, thinkingNoteServiceConfigIns)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", uid_ThinkingNoteServiceConfig), zap.Error(err))
		os.Exit(-1)
	}

	mlog.Logger().Info("> Thinking note service config init finish.")
}

func GetConfig() *thinkingNoteServiceConfig {
	return thinkingNoteServiceConfigIns
}
