package initialize

import (
    "github.com/mats9693/unnamed_plan/services/shared/config"
    "github.com/mats9693/unnamed_plan/services/shared/log"
)

func InitFromConfigCenter(serviceID string, level string, initFunc ...func()) {
    mconfig.InitFromConfigCenter(serviceID, level)
    mlog.Init()

    for i := range initFunc {
        initFunc[i]()
    }
}

func InitFromFile(configFileName string, initFunc ...func()) {
    mconfig.InitFromFile(configFileName)
    mlog.Init()

    for i := range initFunc {
        initFunc[i]()
    }
}
