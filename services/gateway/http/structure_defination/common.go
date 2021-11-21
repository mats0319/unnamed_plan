package structure

import "time"

// request params
const (
	params_UserName   = "userName"
	params_Password   = "password"
	params_OperatorID = "operatorID"
	params_PageSize   = "pageSize"
	params_PageNum    = "pageNum"
	params_Permission = "permission"
	params_UserID     = "userID"
	params_CurrPwd    = "currPwd"
	params_Nickname   = "nickname"

	params_FileName         = "fileName"
	params_ExtensionName    = "extensionName"
	params_LastModifiedTime = "lastModifiedTime"
	params_IsPublic         = "isPublic"
	params_File             = "file"
	params_FileID           = "fileID"

	params_Topic = "topic"
	params_Content = "content"
)

// common

type Total struct {
	Total int `json:"total"`
}

type IsSuccess struct {
	IsSuccess bool `json:"isSuccess"`
}

// user

type UserID struct {
	UserID string `json:"userID"`
}

type Nickname struct {
	Nickname string `json:"nickname"`
}

type Permission struct {
	Permission uint8 `json:"permission"`
}

type Users struct {
	Users []*UserRes `json:"users"`
}

type UserRes struct {
	UserID     string `json:"userID"`
	UserName   string `json:"userName"`
	Nickname   string `json:"nickname"`
	IsLocked   bool   `json:"isLocked"`
	Permission uint8  `json:"permission"`
	CreatedBy  string `json:"createdBy"`
}

// cloud file

type Files struct {
	Files []*FileRes `json:"files"`
}

type FileRes struct {
	FileID           string        `json:"fileID"`
	FileName         string        `json:"fileName"`
	LastModifiedTime time.Duration `json:"lastModifiedTime"`
	FileURL          string        `json:"fileURL"`
	IsPublic         bool          `json:"isPublic"`
	UpdateTime       time.Duration `json:"updateTime"`
	CreatedTime      time.Duration `json:"createdTime"`
}

// note

type Notes struct {
	Notes []*NoteRes `json:"notes"`
}

type NoteRes struct {
	NoteID      string        `json:"noteID"`
	WriteBy     string        `json:"writeBy"`
	Topic       string        `json:"topic"`
	Content     string        `json:"content"`
	IsPublic    bool          `json:"isPublic"`
	UpdateTime  time.Duration `json:"updateTime"`
	CreatedTime time.Duration `json:"createdTime"`
}
