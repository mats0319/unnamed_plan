package initialize

import (
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
)

func InitFromConfigCenter(serviceID string, initFunc ...func() error) error {
	err := mconfig.InitFromConfigCenter(serviceID)
	if err != nil {
		return err
	}

	err = mlog.Init()
	if err != nil {
		return err
	}

	for i := range initFunc {
		err = initFunc[i]()
		if err != nil {
			break
		}
	}

	return err
}

func InitFromFile(configFileName string, initFunc ...func() error) error {
	err := mconfig.InitFromFile(configFileName)
	if err != nil {
		return err
	}

	err = mlog.Init()
	if err != nil {
		return err
	}

	for i := range initFunc {
		err = initFunc[i]()
		if err != nil {
			break
		}
	}

	return err
}

func GetIPAndFreePort() (ip string, port int, err error) {
	ip, err = utils.GetIP()
	if err != nil {
		mlog.Logger().Error("get ip failed", zap.Error(err))
		return
	}

	port, err = utils.GetFreePort()
	if err != nil {
		mlog.Logger().Error("get free port failed", zap.Error(err))
		return
	}

	return
}
