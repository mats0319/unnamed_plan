package model

import "encoding/json"

type ServiceConfig struct {
	ServiceID     string   `pg:",unique:service_config,notnull"`
	Level         string   `pg:",unique:service_config,notnull"`
	ServiceName   string   `pg:",notnull"`
	ConfigItemIDs []string `pg:"config_item_ids"`
	IsDelete      bool     `pg:",use_zero,notnull"`

	Common
}

type ConfigItem struct {
	ConfigItemID   string          `pg:",notnull"`
	ConfigItemName string          `pg:",notnull"`
	ConfigSubItems json.RawMessage `pg:",notnull"`
	BeUsed         []string        // service ids

	Common
}
