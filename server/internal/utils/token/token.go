package token

type Token struct {
	UserName   string    `json:"user_name"`
	Type       TokenType `json:"type"`
	ExpireTime int64     `json:"expire_time"`
}

type TokenType int8

const (
	TokenType_APIAccessToken TokenType = 1
	TokenType_MFAToken                 = 2
)
