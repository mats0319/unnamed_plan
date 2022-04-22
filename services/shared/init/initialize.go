package initialize

import (
    "github.com/mats9693/unnamed_plan/services/shared/config"
    "github.com/mats9693/unnamed_plan/services/shared/log"
)

func Init(configFileName string, initFunc ...func()) {
    mconfig.InitFromFile(configFileName)
    mlog.Init()

    for i := range initFunc {
        initFunc[i]()
    }
}
