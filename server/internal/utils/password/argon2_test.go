package password

import (
	"testing"

	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func TestPassword(t *testing.T) {
	pwdSHA256 := utils.CalcSHA256("123456") // hex(hash('password'))
	t.Log("> Password From Web: ", pwdSHA256)

	pwdArgon2 := GeneratePassword(pwdSHA256)
	t.Log("> Password To DB: ", pwdArgon2)

	err := VerifyPassword(pwdSHA256, pwdArgon2)
	if err != nil {
		t.Fatal(err.String())
	}
	t.Log("> Verified.")
}
