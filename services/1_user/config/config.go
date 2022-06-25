package config

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
)

type userServiceConfig struct {
	init bool

	ARankAdminPermission uint8  `json:"ARankAdminPermission"`
	SRankAdminPermission uint8  `json:"SRankAdminPermission"`
}

var userServiceConfigIns = &userServiceConfig{}

func Init() error {
	if userServiceConfigIns.init { // have initialized
		mlog.Logger().Error("already initialized")
		return nil
	}

	byteSlice := mconfig.GetConfig(mconst.UID_Service_User)

	err := json.Unmarshal(byteSlice, userServiceConfigIns)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", mconst.UID_Service_User), zap.Error(err))
		return utils.NewError(err.Error())
	}

	userServiceConfigIns.init = true

	mlog.Logger().Info("> User service config init finish.")

	return nil
}

func GetConfig() *userServiceConfig {
	return userServiceConfigIns
}
