package main

import (
	"encoding/json"
	"github.com/go-pg/pg/v10"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"log"
)

var serviceConfig = []*model.ServiceConfig{
	// todo: impl
	{
		ServiceID:   "84d1fecc-3be9-439e-8144-209ffc00a975",
		Level:       mconst.ConfigLevel_Dev,
		ServiceName: "gateway",
		ConfigItemIDs: []string{
			"1cd10cb8-ecf5-4855-a886-76b148ed104a",
			"3b839c1f-9f1e-474b-b3da-7b5e9bc792ec",
		},
		IsDelete: false,
		Common:   model.NewCommon(),
	}, // gateway service, dev
	{
		ServiceID:   "84d1fecc-3be9-439e-8144-209ffc00a975",
		Level:       mconst.ConfigLevel_Production,
		ServiceName: "gateway",
		ConfigItemIDs: []string{
			"1cd10cb8-ecf5-4855-a886-76b148ed104a",
			"3b839c1f-9f1e-474b-b3da-7b5e9bc792ec",
		},
		IsDelete: false,
		Common:   model.NewCommon(),
	}, // gateway service, production
}

var configItem = []*model.ConfigItem{
	// todo: impl
	{
		ConfigItemID:   "1cd10cb8-ecf5-4855-a886-76b148ed104a",
		ConfigItemName: "gateway rpc client",
		ConfigSubItems: json.RawMessage(`{
		  "userClientTarget": "127.0.0.1:10001",
		  "cloudFileClientTarget": "127.0.0.1:10002",
		  "noteClientTarget": "127.0.0.1:10003",
		  "taskClientTarget": "127.0.0.1:10004"
		}`),
		BeUsed: []string{
			"84d1fecc-3be9-439e-8144-209ffc00a975",
		},
		Common: model.NewCommon(),
	},
	{
		ConfigItemID:   "3b839c1f-9f1e-474b-b3da-7b5e9bc792ec",
		ConfigItemName: "gateway http",
		ConfigSubItems: json.RawMessage(`{
		  "port": "9693",
		  "sources": ["public web", "admin web"],
		  "unlimitedSources": ["public mobile"],
		  "limitMultiLoginConfig": {
		    "limitMultiLogin": true,
		    "keepTokenValid": 7200
		  }
		}`),
		BeUsed: []string{
			"84d1fecc-3be9-439e-8144-209ffc00a975",
		},
		Common: model.NewCommon(),
	},
}

func setDefaultServiceConfig(db *pg.DB) {
	for i := range serviceConfig {
		_, err := db.Model(serviceConfig[i]).Insert()
		if err != nil {
			log.Printf("set default service config failed, index: %d, error: %v\n", i, err)
		}
	}
}

func setDefaultConfigItem(db *pg.DB) {
	for i := range configItem {
		_, err := db.Model(configItem[i]).Insert()
		if err != nil {
			log.Printf("set default config item failed, index: %d, error: %v\n", i, err)
		}
	}
}
