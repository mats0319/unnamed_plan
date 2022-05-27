package dao

import "github.com/mats9693/unnamed_plan/services/shared/db/model"

type ServiceConfigDao interface {
	// Query get all records except 'is delete' ones
	Query() ([]*model.ServiceConfig, error)
}

type ConfigItemDao interface {
	// Query get all records
	Query() ([]*model.ConfigItem, error)
}
