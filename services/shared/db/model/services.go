package model

import (
	"time"
)

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

type ThinkingNote struct {
	NoteID    string `pg:",unique,notnull"`
	WriteBy   string `pg:",notnull"` // user id
	Topic     string
	Content   string `pg:",notnull"`
	IsPublic  bool   `pg:",use_zero,notnull"`
	IsDeleted bool   `pg:",use_zero,notnull"`

	Common
}
