package structure

import (
	mconst "github.com/mats9693/unnamed_plan/services/shared/const"
	"time"
)

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

	params_Topic   = "topic"
	params_Content = "content"
	params_NoteID = "noteID"

	params_TaskName = "taskName"
	params_Description = "description"
	params_PreTaskIDs = "preTaskIDs"
	params_TaskID = "taskID"
	params_Status = "status"
)

// common

type Total struct {
	Total uint32 `json:"total"`
}

// user

type UserID struct {
	UserID string `json:"userID"`
}

type Nickname struct {
	Nickname string `json:"nickname"`
}

type Permission struct {
	Permission uint32 `json:"permission"`
}

type Users struct {
	Users []*UserRes `json:"users"`
}

type UserRes struct {
	UserID     string `json:"userID"`
	UserName   string `json:"userName"`
	Nickname   string `json:"nickname"`
	IsLocked   bool   `json:"isLocked"`
	Permission uint32 `json:"permission"`
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

// task

type Tasks struct {
	Tasks []*TaskRes `json:"tasks"`
}

type TaskRes struct {
	TaskID      string            `json:"taskID"`
	TaskName    string            `json:"taskName"`
	Description string            `json:"description"`
	PreTaskIDs  []string          `json:"preTaskIDs"`
	Status      mconst.TaskStatus `json:"status"`
	UpdateTime  time.Duration     `json:"updateTime"`
	CreatedTime time.Duration     `json:"createdTime"`
}
