package password

import (
	"crypto/subtle"

	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func GeneratePassword(pwdSHA256 string) (pwdArgon2 string) {
	params := defaultAlgorithmParams()

	key := params.deriveKey(pwdSHA256)
	pwdArgon2 = params.encode(key)

	return
}

func VerifyPassword(pwdSHA256 string, pwdArgon2 string) *utils.Error {
	params, oldKey, e := (&AlgorithmParams{}).decode(pwdArgon2)
	if e != nil {
		return e
	}

	newKey := params.deriveKey(pwdSHA256)

	if subtle.ConstantTimeCompare(oldKey, newKey) != 1 { // 使用恒定时间比较防止时序攻击
		// 这里不打印错误，因为部分应用场景要求密码验证不能通过（例如修改密码时，新、旧密码不能一样）。
		// 具体的，假设这里打印错误，那么在修改密码时，即使一切正确执行，控制台也会报错。
		// 虽然这样可能导致前面decode出错时，错误被打印2次，但我认为这是可以接受的
		return utils.ErrWrongPassword().WithParam("old key", oldKey).WithParam("new key", newKey)
	}

	return nil
}
