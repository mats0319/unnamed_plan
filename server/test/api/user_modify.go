package api

import (
	"fmt"

	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ModifyUser() {
	testCase("no changes", modifyUserCase_NoChanges)
	testCase("same pwd", modifyUserCase_SamePwd)
	testCase("success", modifyUserCase_Success)
}

func modifyUserParams(nickname string, password string) string {
	return fmt.Sprintf(`{"nickname":"%s","password":"%s"}`, nickname, password)
}

func modifyUserCase_NoChanges() string {
	res := httpInvoke(api.URI_ModifyUser, modifyUserParams("", ""), accessToken_User)
	if res.IsSuccess || !errorIs(res.Err, utils.ErrNoChanges()) {
		return unknownError
	}

	return ""
}

func modifyUserCase_SamePwd() string {
	res := httpInvoke(api.URI_ModifyUser, modifyUserParams("", pwdSHA256), accessToken_User)
	if res.IsSuccess || !errorIs(res.Err, utils.ErrSamePassword()) {
		return unknownError
	}

	return ""
}

func modifyUserCase_Success() string {
	res := httpInvoke(api.URI_ModifyUser, modifyUserParams("new nickname", ""), accessToken_User)
	if !res.IsSuccess {
		return res.Err
	}

	user, err := dal.GetUser("user")
	if (user != nil && user.Nickname != "new nickname") || err != nil {
		return unknownError
	}

	return ""
}
