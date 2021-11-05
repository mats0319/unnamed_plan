package model

import (
	"github.com/mats9693/utils/uuid"
	"time"
)

type Common struct {
	ID          string        `pg:",pk"`
	CreatedTime time.Duration `pg:",use_zero,notnull"`
	UpdateTime  time.Duration `pg:",use_zero,notnull"`
	// todo: 乐观锁
}

type User struct {
	UserID     string `pg:",unique,notnull"`
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

func NewCommon() Common {
	currTime := time.Duration(time.Now().Unix())

	return Common{
		ID:          uuid.New(),
		UpdateTime:  currTime,
		CreatedTime: currTime,
	}
}
