package model

import (
	"time"
	"github.com/mats9693/utils/uuid"
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
	FileID           string `pg:",unique"` // sha256('user id' + timestamp), storage name
	UploadedBy       string // user id
	FileName         string // show name, for display
	ExtensionName    string
	LastModifiedTime time.Duration `pg:",use_zero"`
	FileSize         int64         `pg:",use_zero"`
	IsPublic         bool          `pg:",use_zero"`
	IsDeleted        bool          `pg:",use_zero"`

	Common
}

type ThinkingNote struct {
	NoteID    string `pg:",unique"`
	WriteBy   string // user id
	Topic     string
	Content   string `pg:",notnull"`
	IsPublic  bool   `pg:",use_zero"`
	IsDeleted bool   `pg:",use_zero"`

	Common
}

type GameResult struct {
	ResultID string        `pg:",unique"`
	Player   string        // user id
	Duration time.Duration `pg:",use_zero"` // unit: second
	Result   string        // json string
	Remark   string        // for extend

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
