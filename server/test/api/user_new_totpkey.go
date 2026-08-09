package api

import api "github.com/mats0319/unnamed_plan/server/cmd/api/go"

var totpKey = ""

func NewTOTPKey() {
	testCase("success", newTOTPKeyCase_Success)
}

func newTOTPKeyCase_Success() string {
	res := httpInvoke(api.URI_NewTOTPKey, ``, accessToken_User)
	if !res.IsSuccess {
		return res.Err
	}

	if len(res.Data.TOTPKey) != 16 {
		return unknownError
	}

	totpKey = res.Data.TOTPKey

	return ""
}
