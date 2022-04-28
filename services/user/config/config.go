package config

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
	"os"
)

type userServiceConfig struct {
	init bool

	Address              string `json:"address"`
	ARankAdminPermission uint8  `json:"ARankAdminPermission"`
	SRankAdminPermission uint8  `json:"SRankAdminPermission"`
}

var userServiceConfigIns = &userServiceConfig{}

func Init() {
	if userServiceConfigIns.init { // have initialized
		return
	}

	byteSlice := mconfig.GetConfig(mconst.UID_Service_User)

	err := json.Unmarshal(byteSlice, userServiceConfigIns)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", mconst.UID_Service_User), zap.Error(err))
		os.Exit(-1)
	}

	userServiceConfigIns.init = true

	mlog.Logger().Info("> User service config init finish.")
}

func GetConfig() *userServiceConfig {
	return userServiceConfigIns
}
