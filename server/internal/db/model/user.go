package model

type User struct {
	Model
	UserName string `gorm:"unique;not null"` // login name, can't modify
	Nickname string `gorm:"not null"`        // display name
	Password string `gorm:"not null"`
	IsAdmin  bool

	EnableMFA bool
	TOTPKey   string

	// last login, use 'user.UpdateAt'
}
