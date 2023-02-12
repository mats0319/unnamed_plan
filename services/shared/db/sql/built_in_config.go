package main

import (
	"encoding/json"
	"github.com/go-pg/pg/v10"
	mconst "github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"log"
)

var serviceConfig = []*model.ServiceConfig{
	{
		ServiceID:   "84d1fecc-3be9-439e-8144-209ffc00a975",
		Level:       mconst.ConfigLevel_Default,
		ServiceName: "gateway",
		ConfigItemIDs: []string{
			configItem[0].Common.ID,
			configItem[1].Common.ID,
		},
		IsDelete: false,
		Common:   model.NewCommon(),
	}, // gateway service
	{
		ServiceID:   "eafbda7d-c951-4fc9-8b45-8c90189c1e36",
		Level:       mconst.ConfigLevel_Dev,
		ServiceName: "user",
		ConfigItemIDs: []string{
			configItem[2].Common.ID,
			configItem[4].Common.ID,
		},
		IsDelete: false,
		Common:   model.NewCommon(),
	}, // user service, dev
	{
		ServiceID:   "eafbda7d-c951-4fc9-8b45-8c90189c1e36",
		Level:       mconst.ConfigLevel_Production,
		ServiceName: "user",
		ConfigItemIDs: []string{
			configItem[3].Common.ID,
			configItem[4].Common.ID,
		},
		IsDelete: false,
		Common:   model.NewCommon(),
	}, // user service, production
	{
		ServiceID:   "1b5ab1d2-de6d-4377-9a4e-a184b24d1a0f",
		Level:       mconst.ConfigLevel_Dev,
		ServiceName: "cloud file",
		ConfigItemIDs: []string{
			configItem[2].Common.ID,
			configItem[5].Common.ID,
		},
		IsDelete: false,
		Common:   model.NewCommon(),
	}, // cloud file service, dev
	{
		ServiceID:   "1b5ab1d2-de6d-4377-9a4e-a184b24d1a0f",
		Level:       mconst.ConfigLevel_Production,
		ServiceName: "cloud file",
		ConfigItemIDs: []string{
			configItem[3].Common.ID,
			configItem[6].Common.ID,
		},
		IsDelete: false,
		Common:   model.NewCommon(),
	}, // cloud file service, production
}

var configItem = []*model.ConfigItem{
	{
		ConfigItemID:   "3b839c1f-9f1e-474b-b3da-7b5e9bc792ec",
		ConfigItemName: "gateway http",
		ConfigSubItems: json.RawMessage(`{
		  "port": "9693",
		  "sources": ["web"]
		}`),
		Common: model.NewCommon(),
	}, // gateway http
	{
		ConfigItemID:   "5c581355-06cf-465e-ae63-5b9c0284477d",
		ConfigItemName: "gateway plugins",
		ConfigSubItems: json.RawMessage(`{
		  "limitMultiLoginConfig": {
		    "limitMultiLogin": true,
			"sources": ["web"],
		    "keepTokenValid": 7200
		  }
		}`),
		Common: model.NewCommon(),
	}, // gateway plugins
	{
		ConfigItemID:   "658e06f7-71d5-4ada-b715-8c1a4489e5d2",
		ConfigItemName: "db _ " + mconst.ConfigLevel_Dev,
		ConfigSubItems: json.RawMessage(`{
          "DBMS": "postgresql",
          "addr": "117.50.177.201:5432",
          "dbName": "unnamed_plan_test",
          "user": "mario_test",
          "password": "123456",
          "timeout": 3,
          "showSQL": true
        }`),
		Common: model.NewCommon(),
	}, // db, dev
	{
		ConfigItemID:   "658e06f7-71d5-4ada-b715-8c1a4489e5d2",
		ConfigItemName: "db _ " + mconst.ConfigLevel_Production,
		ConfigSubItems: json.RawMessage(`{
          "DBMS": "postgresql",
          "addr": "127.0.0.1:5432",
          "dbName": "unnamed_plan_test",
          "user": "mario_test",
          "password": "123456",
          "timeout": 3,
          "showSQL": true
        }`),
		Common: model.NewCommon(),
	}, // db, production
	{
		ConfigItemID:   "eafbda7d-c951-4fc9-8b45-8c90189c1e36",
		ConfigItemName: "user service basic config",
		ConfigSubItems: json.RawMessage(`{
          "ARankAdminPermission": 6,
          "SRankAdminPermission": 8
        }`),
		Common: model.NewCommon(),
	}, // user service basic config
	{
		ConfigItemID:   "1b5ab1d2-de6d-4377-9a4e-a184b24d1a0f",
		ConfigItemName: "cloud file service basic config _ " + mconst.ConfigLevel_Dev,
		ConfigSubItems: json.RawMessage(`{
          "cloudFileRootPath": "",
          "cloudFilePublicDir": "public"
        }`),
		Common: model.NewCommon(),
	}, // cloud file service basic config, dev
	{
		ConfigItemID:   "1b5ab1d2-de6d-4377-9a4e-a184b24d1a0f",
		ConfigItemName: "cloud file service basic config _ " + mconst.ConfigLevel_Production,
		ConfigSubItems: json.RawMessage(`{
          "cloudFileRootPath": "/home/mario/cloud_file/",
          "cloudFilePublicDir": "public"
        }`),
		Common: model.NewCommon(),
	}, // cloud file service basic config, production
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
