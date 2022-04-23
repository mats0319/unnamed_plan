package model

import "encoding/json"

type ServiceConfig struct {
	ServiceID     string   `pg:",unique:service_config,notnull"`
	Level         string   `pg:",unique:service_config,notnull"`
	ServiceName   string   `pg:",notnull"`
	ConfigItemIDs []string `pg:"config_item_ids"`
	IsDelete      bool     `pg:",use_zero,notnull"`
}

type ConfigItem struct {
	ConfigItemID   string `pg:",notnull"`
	ConfigItemName string `pg:",notnull"`
	ConfigItemTag  string
	ConfigSubItems json.RawMessage `pg:",notnull"`
	UsedIn         []string        // service ids
}
