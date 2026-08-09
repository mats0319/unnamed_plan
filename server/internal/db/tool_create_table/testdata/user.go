package testdata

import (
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/password"
)

var TestUsers = []*model.User{
	NewUser("admin", "", true, true, "5SSFNNEJUENPCCKP"),
	NewUser("user", "", false, false, ""),
}

func NewUser(userName string, nickname string, isAdmin bool, enable2FA bool, totpKey string) *model.User {
	pwdSHA256 := utils.CalcSHA256("123456")

	if len(nickname) < 1 {
		nickname = userName
	}

	keyEncrypted, _ := utils.Encrypt(totpKey, "cBsnYH1yDvFfW4q84wAz7FhrzWHiUjQk")

	return &model.User{
		UserName:  userName,
		Nickname:  nickname,
		Password:  password.GeneratePassword(pwdSHA256),
		IsAdmin:   isAdmin,
		EnableMFA: enable2FA,
		TOTPKey:   keyEncrypted,
	}
}
