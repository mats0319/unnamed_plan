package main

import (
    "github.com/mats9693/unnamed_plan/services/config_center/db"
    i "github.com/mats9693/unnamed_plan/services/config_center/init"
    "github.com/mats9693/unnamed_plan/services/shared/db"
    "github.com/mats9693/unnamed_plan/services/shared/init"
)

func main() {
    initialize.InitFromFile("config.json", mdb.Init, db.Init, i.InitSupportConfig)
}
