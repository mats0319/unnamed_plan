package api

import (
	"fmt"

	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func Register() {
	testCase("duplicate register", registerCase_Duplicate)
	testCase("success", registerCase_Success)
}

func registerParams(userName string, password string) string {
	return fmt.Sprintf(`{"user_name":"%s","password":"%s"}`, userName, password)
}

func registerCase_Duplicate() string {
	res := httpInvoke(api.URI_Register, registerParams("admin", "123"), "")
	if res.IsSuccess || !errorIs(res.Err, utils.ErrUserExist()) {
		return unknownError
	}

	return ""
}

func registerCase_Success() string {
	res := httpInvoke(api.URI_Register, registerParams("new_user", pwdSHA256), "")
	if !res.IsSuccess {
		return res.Err
	}

	count, data, err := dal.ListUsers(10, 1)
	if count != 4 || len(data) != 4 || err != nil {
		return unknownError
	}

	return ""
}
