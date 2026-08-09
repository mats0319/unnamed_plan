package api

import (
	"fmt"

	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
)

func SetMFAStatus() {
	testCase("success - enable MFA", setMFAStatusCase_SuccessEnable)
	testCase("success - disable MFA", setMFAStatusCase_SuccessDisable)
}

func setMFAStatusParams(enableMFA bool, applyNewKeyFlag bool, totpKey string) string {
	return fmt.Sprintf(`{"enable_mfa":%t,"apply_new_key_flag":%t,"totp_code":"%s"}`,
		enableMFA, applyNewKeyFlag, calcTotpCode(totpKey))
}

func setMFAStatusCase_SuccessEnable() string {
	res := httpInvoke(api.URI_SetMFAStatus, setMFAStatusParams(true, true, totpKey), accessToken_User)
	if !res.IsSuccess {
		return res.Err
	}

	user, e := dal.GetUser("user")
	// 这里判断totp key的长度，因为计划添加totp key加密保存功能，那时就没有办法比较具体值了，所以这里直接判断不为空
	if e != nil || !user.EnableMFA || len(user.TOTPKey) < 1 {
		return unknownError
	}

	return ""
}

func setMFAStatusCase_SuccessDisable() string {
	res := httpInvoke(api.URI_SetMFAStatus, setMFAStatusParams(false, false, ""), accessToken_User)
	if !res.IsSuccess {
		return res.Err
	}

	user, e := dal.GetUser("user")
	if e != nil || user.EnableMFA {
		return unknownError
	}

	return ""
}
