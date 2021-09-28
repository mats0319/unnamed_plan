package http_res_type

import "time"

type HTTPResFiles struct {
	FileID      string        `json:"fileID"`
	FileName    string        `json:"fileName"`
	FileURL     string        `json:"fileURL"`
	IsPublic    bool          `json:"isPublic"`
	UpdateTime  time.Duration `json:"updateTime"`
	CreatedTime time.Duration `json:"createdTime"`
}

type HttpResUser struct {
	UserID     string `json:"userID"`
	Nickname   string `json:"nickname"`
	IsLocked   bool   `json:"isLocked"`
	Permission uint8  `json:"permission"`
	CreatedBy  string `json:"createdBy"`
}
