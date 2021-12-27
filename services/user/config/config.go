package config

import (
	"encoding/json"
	"github.com/mats9693/utils/toy_server/config"
	"github.com/mats9693/utils/toy_server/log"
	"go.uber.org/zap"
	"os"
)

const uid_UserServiceConfig = "6d31fac1-346c-4b03-8596-8b3e1ee5b960"

type userServiceConfig struct {
	Address              string `json:"address"`
	ARankAdminPermission uint8  `json:"ARankAdminPermission"`
	SRankAdminPermission uint8  `json:"SRankAdminPermission"`
}

var userServiceConfigIns = &userServiceConfig{}

func init() {
	byteSlice := mconfig.GetConfig(uid_UserServiceConfig)

	err := json.Unmarshal(byteSlice, userServiceConfigIns)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", uid_UserServiceConfig), zap.Error(err))
		os.Exit(-1)
	}

	mlog.Logger().Info("> User service config init finish.")
}

func GetConfig() *userServiceConfig {
	return userServiceConfigIns
}
