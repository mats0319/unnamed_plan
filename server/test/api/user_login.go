package api

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"math"
	"time"

	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func Login() {
	testCase("user not exist", loginCase_UserNotExist)
	testCase("wrong pwd", loginCase_WrongPwd)
	testCase("success - disable 2fa", loginCase_SuccessDisable2fa)
	testCase("wrong totp code", loginCase_WrongTOTPCode)
	testCase("success - enable 2fa", loginCase_SuccessEnable2fa)
}

func loginParams(userName string, password string) string {
	return fmt.Sprintf(`{"user_name":"%s","password":"%s"}`, userName, password)
}

func loginMFAParams(mfaToken string, totpKey string) string {
	return fmt.Sprintf(`{"mfa_token":"%s","totp_code":"%s"}`, mfaToken, calcTotpCode(totpKey))
}

func loginCase_UserNotExist() string {
	res := httpInvoke(api.URI_Login, loginParams("not exist", "123"), "")
	if res.IsSuccess || !errorIs(res.Err, utils.ErrUserNotFound()) {
		return unknownError
	}

	return ""
}

func loginCase_WrongPwd() string {
	res := httpInvoke(api.URI_Login, loginParams("admin", "wrong pwd"), "")
	if res.IsSuccess || !errorIs(res.Err, utils.ErrWrongPassword()) {
		return unknownError
	}

	return ""
}

func loginCase_SuccessDisable2fa() string {
	res := httpInvoke(api.URI_Login, loginParams("user", pwdSHA256), "")
	if !res.IsSuccess {
		return res.Err
	}

	return ""
}

func loginCase_WrongTOTPCode() string {
	res := httpInvoke(api.URI_Login, loginParams("user_with_totp", pwdSHA256), "")
	if !res.IsSuccess {
		return res.Err
	}

	mfaToken := res.Data.MFAToken
	if len(mfaToken) < 1 {
		return unknownError
	}

	res = httpInvoke(api.URI_LoginMFA, loginMFAParams(mfaToken, ""), "")
	if res.IsSuccess || !errorIs(res.Err, utils.ErrWrongTOTPCode()) {
		return unknownError
	}

	return ""
}

func loginCase_SuccessEnable2fa() string {
	res := httpInvoke(api.URI_Login, loginParams("user_with_totp", pwdSHA256), "")
	if !res.IsSuccess {
		return res.Err
	}

	mfaToken := res.Data.MFAToken
	if len(mfaToken) < 1 {
		return unknownError
	}

	res = httpInvoke(api.URI_LoginMFA, loginMFAParams(mfaToken, "NVQXE2LP"), "")
	if !res.IsSuccess {
		return res.Err
	}

	return ""
}

func calcTotpCode(keyBase32 string) string {
	if keyBase32 == "" {
		return "000000"
	}

	key := make([]byte, 10)
	n, err := base32.StdEncoding.Decode(key, []byte(keyBase32))
	if err != nil {
		return ""
	}
	key = key[:n]

	hasher := hmac.New(sha1.New, key)
	hasher.Write(iTob(time.Now().Unix() / 30))
	hmacHash := hasher.Sum(nil)

	offset := int(hmacHash[len(hmacHash)-1] & 0x0f)
	// 算法要求屏蔽最高有效位
	longPassword := int(hmacHash[offset]&0x7f)<<24 |
		int(hmacHash[offset+1])<<16 |
		int(hmacHash[offset+2])<<8 |
		int(hmacHash[offset+3])

	password := longPassword % int(math.Pow10(6))

	return fmt.Sprintf("%06d", password)
}

func iTob(integer int64) []byte {
	byteSlice := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		byteSlice[i] = byte(integer & 0xff)
		integer = integer >> 8
	}
	return byteSlice
}
