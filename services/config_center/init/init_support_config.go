package initialize

import (
	"github.com/mats9693/unnamed_plan/services/config_center/db"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
	"os"
)

var (
	serviceConfigMap = make(map[string]*configWrapper) // service id - config wrapper

	serviceConfigs []*model.ServiceConfig
	configItems    []*model.ConfigItem
)

type configWrapper struct {
	Default    *mconfig.Config
	Dev        *mconfig.Config
	Production *mconfig.Config
	Test       *mconfig.Config
}

// GetServiceConfig return nil if un-support
func GetServiceConfig(serviceID string, level string) (res *mconfig.Config) {
	config := serviceConfigMap[serviceID]

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

func InitSupportConfig() {
	var err error
	serviceConfigs, err = db.GetServiceConfigDao().Query()
	if err != nil {
		mlog.Logger().Error("get service config failed", zap.Error(err))
		os.Exit(-1)
	}

	configItems, err = db.GetConfigItemDao().Query()
	if err != nil {
		mlog.Logger().Error("get config item failed", zap.Error(err))
		os.Exit(-1)
	}

	supportValidServiceConfig()
}

func supportValidServiceConfig() {
	for i := range serviceConfigs {
		serviceConfigItem := serviceConfigs[i]

		// a 'service config' is valid when all 'config item ids' are existed
		config := &mconfig.Config{Level: serviceConfigItem.Level}

		// if 'service config' is valid, format and support it
		configItemList := make([]*mconfig.ConfigItem, 0, len(serviceConfigItem.ConfigItemIDs))
		isValid := false
		for j := range serviceConfigItem.ConfigItemIDs {
			isValid = false

			configItemID := serviceConfigItem.ConfigItemIDs[j]
			for k := range configItems {
				if configItemID == configItems[k].ID {
					item := configItems[k]

					configItemList = append(configItemList, &mconfig.ConfigItem{
						UID:  item.ConfigItemID,
						Name: item.ConfigItemName,
						Json: item.ConfigSubItems,
					})

					isValid = true
				}
			}

			if !isValid {
				// under normal circumstances, all service config record are valid, if not, log it
				mlog.Logger().Error("un-support service config",
					zap.String("service id", serviceConfigItem.ID),
					zap.String("config item id", configItemID))
				break
			}
		}

		if !isValid {
			continue
		}

		// support valid 'service config'
		config.ConfigItems = configItemList

		serviceConfig, ok := serviceConfigMap[serviceConfigItem.ServiceID]
		if !ok {
			serviceConfig = &configWrapper{}
		}

		switch serviceConfigItem.Level {
		case mconst.ConfigLevel_Default:
			serviceConfig.Default = config
		case mconst.ConfigLevel_Dev:
			serviceConfig.Dev = config
		case mconst.ConfigLevel_Production:
			serviceConfig.Production = config
		case mconst.ConfigLevel_Test:
			serviceConfig.Test = config
		}

		serviceConfigMap[serviceConfigItem.ServiceID] = serviceConfig
	}
}
