package initialize

import (
	"github.com/mats9693/unnamed_plan/services/core/db"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
	"sync"
)

var serviceConfigMap = sync.Map{} // service id - *configWrapper

type configWrapper struct {
	Default    *mconfig.Config
	Dev        *mconfig.Config
	Production *mconfig.Config
	Test       *mconfig.Config
}

// GetServiceConfig return nil if un-support
func GetServiceConfig(serviceID string, level string) (res *mconfig.Config) {
	configI, ok := serviceConfigMap.Load(serviceID)
	if !ok {
		mlog.Logger().Info("unsupported service id: " + serviceID)
		return
	}

	config, _ := configI.(*configWrapper)

	switch level {
	case mconst.ConfigLevel_Default:
		res = config.Default
	case mconst.ConfigLevel_Dev:
		res = config.Dev
	case mconst.ConfigLevel_Production:
		res = config.Production
	case mconst.ConfigLevel_Test:
		res = config.Test
	}

	// if special level is un-support, try 'default' level
	if res == nil {
		res = config.Default
	}

	return
}

func InitSupportedConfig() error {
	serviceConfigSlice, err := db.GetServiceConfigDao().Query()
	if err != nil {
		mlog.Logger().Error("get service config failed", zap.Error(err))
		return err
	}

	configItemSlice, err := db.GetConfigItemDao().Query()
	if err != nil {
		mlog.Logger().Error("get config item failed", zap.Error(err))
		return err
	}

	supportValidServiceConfig(serviceConfigSlice, configItemSlice)

	return nil
}

func supportValidServiceConfig(serviceConfigSlice []*model.ServiceConfig, configItemSlice []*model.ConfigItem) {
	for i := range serviceConfigSlice {
		serviceConfigRecord := serviceConfigSlice[i]

		serviceConfigIns := &mconfig.Config{Level: serviceConfigRecord.Level}

		// a 'service config' is valid means all its 'config item ids' are existed
		validConfigItems := make([]*mconfig.ConfigItem, 0, len(serviceConfigRecord.ConfigItemIDs))
		isValid := true
		for j := 0; j < len(serviceConfigRecord.ConfigItemIDs) && isValid; j++ {
			isValid = false

			index := getIndex(configItemSlice, serviceConfigRecord.ConfigItemIDs[j])
			if index >= 0 {
				isValid = true

				validConfigItems = append(validConfigItems, formatConfigItem(configItemSlice[index]))
			}
		}

		if !isValid {
			continue
		}

		serviceConfigIns.ConfigItems = validConfigItems

		// support valid 'service config'
		var serviceConfigWrapper *configWrapper
		{
			serviceConfigI, ok := serviceConfigMap.Load(serviceConfigRecord.ServiceID)
			if !ok { // first config of target service
				serviceConfigWrapper = &configWrapper{}
			} else {
				serviceConfigWrapper, _ = serviceConfigI.(*configWrapper)
			}
		}

		switch serviceConfigRecord.Level {
		case mconst.ConfigLevel_Default:
			serviceConfigWrapper.Default = serviceConfigIns
		case mconst.ConfigLevel_Dev:
			serviceConfigWrapper.Dev = serviceConfigIns
		case mconst.ConfigLevel_Production:
			serviceConfigWrapper.Production = serviceConfigIns
		case mconst.ConfigLevel_Test:
			serviceConfigWrapper.Test = serviceConfigIns
		}

		serviceConfigMap.Store(serviceConfigRecord.ServiceID, serviceConfigWrapper)
	}
}

func getIndex(array []*model.ConfigItem, configItemID string) int {
	index := -1
	for i := range array {
		if configItemID == array[i].ID {
			index = i
			break
		}
	}

	return index
}

func formatConfigItem(data *model.ConfigItem) *mconfig.ConfigItem {
	if data == nil {
		return &mconfig.ConfigItem{}
	}

	return &mconfig.ConfigItem{
		UID:  data.ConfigItemID,
		Name: data.ConfigItemName,
		Json: data.ConfigSubItems,
	}
}
