package model

import (
	"encoding/json"
	"github.com/mats9693/utils/uuid"
	"time"
)

// config center

type ServiceConfig struct {
	ServiceID     string   `pg:",unique:service_config,notnull"`
	Level         string   `pg:",unique:service_config,notnull"`
	ServiceName   string   `pg:",notnull"`
	ConfigItemIDs []string `pg:"config_item_ids"` // config item db record ids
	IsDelete      bool     `pg:",use_zero,notnull"`

	Common
}

type ConfigItem struct {
	ConfigItemID   string          `pg:",notnull"`
	ConfigItemName string          `pg:",notnull"`
	ConfigSubItems json.RawMessage `pg:",notnull"`

	Common
}

// services

type User struct {
	UserName   string `pg:",unique,notnull"` // login name
	Nickname   string `pg:",notnull"`        // show name
	Password   string `pg:"type:varchar(64),notnull"`
	Salt       string `pg:"type:varchar(10),notnull"`
	IsLocked   bool   `pg:",use_zero,notnull"`
	Permission uint8  `pg:",use_zero,notnull"`
	CreatedBy  string `pg:",notnull"`

	Common
}

type CloudFile struct {
	FileID           string        `pg:",unique,notnull"` // sha256('user id' + timestamp), storage name
	UploadedBy       string        `pg:",notnull"`        // user id
	FileName         string        `pg:",notnull"`        // display name
	ExtensionName    string        `pg:",notnull"`
	LastModifiedTime time.Duration `pg:",use_zero,notnull"`
	FileSize         int64         `pg:",use_zero,notnull"`
	IsPublic         bool          `pg:",use_zero,notnull"`
	IsDeleted        bool          `pg:",use_zero,notnull"`

	Common
}

// common

type Common struct {
	ID             string        `pg:",pk"`
	CreatedTime    time.Duration `pg:",use_zero,notnull"`
	UpdateTime     time.Duration `pg:",use_zero,notnull"`
	OptimisticLock uint64        `pg:",use_zero,notnull"`
}

func NewCommon() Common {
	currTime := time.Duration(time.Now().Unix())

	return Common{
		ID:             uuid.New(),
		UpdateTime:     currTime,
		CreatedTime:    currTime,
		OptimisticLock: 0,
	}
}
