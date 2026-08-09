package testdata

import (
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/password"
)

func PresetUser() []*model.User {
	pwdSHA256 := utils.CalcSHA256("123456")
	totpKey, _ := utils.Encrypt("NVQXE2LP", "cBsnYH1yDvFfW4q84wAz7FhrzWHiUjQk")

	return []*model.User{
		{
			UserName: "admin",
			Nickname: "admin",
			Password: password.GeneratePassword(pwdSHA256),
			IsAdmin:  true,
		},
		{
			UserName: "user",
			Nickname: "user",
			Password: password.GeneratePassword(pwdSHA256),
		},
		{
			UserName:  "user_with_totp",
			Nickname:  "user_with_totp",
			Password:  password.GeneratePassword(pwdSHA256),
			EnableMFA: true,
			TOTPKey:   totpKey, // base32 of 'mario'
		},
	}
}
