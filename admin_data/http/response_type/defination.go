package http_res_type

import "time"

type HttpResUser struct {
	UserID     string `json:"userID"`
	UserName   string `json:"userName"`
	Nickname   string `json:"nickname"`
	IsLocked   bool   `json:"isLocked"`
	Permission uint8  `json:"permission"`
	CreatedBy  string `json:"createdBy"`
}

type HTTPResFiles struct {
	FileID      string        `json:"fileID"`
	FileName    string        `json:"fileName"`
	FileURL     string        `json:"fileURL"`
	IsPublic    bool          `json:"isPublic"`
	UpdateTime  time.Duration `json:"updateTime"`
	CreatedTime time.Duration `json:"createdTime"`
}
