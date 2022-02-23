package mresponse

type ResponseData struct {
	UserID string `json:"userID,omitempty"`
	Token  string `json:"token,omitempty"`

	HasError bool   `json:"hasError"`
	Data     string `json:"data"`
}
