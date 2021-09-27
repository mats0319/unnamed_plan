package model

import (
	"github.com/mats9693/utils/uuid"
	"time"
)

type Common struct {
	ID          string        `pg:",pk"`
	CreatedTime time.Duration `pg:",use_zero"`
	UpdateTime  time.Duration `pg:",use_zero"`
	// todo: 乐观锁
}

type User struct {
	UserID     string `pg:",unique"`
	UserName   string `pg:",unique"`  // login name
	Nickname   string `pg:",notnull"` // show name
	Password   string `pg:"type:varchar(64),notnull"`
	Salt       string `pg:",notnull"`
	IsLocked   bool   `pg:",use_zero"`
	Permission uint8  `pg:",use_zero"`
	CreatedBy  string

	Common
}

type CloudFile struct {
	FileID        string `pg:",unique"` // sha256('user id' + timestamp), storage name
	UploadedBy    string // user id
	FileName      string // show name, for display
	ExtensionName string
	FileSize      string
	IsPublic      bool `pg:",use_zero"`
	IsDeleted     bool `pg:",use_zero"`

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
