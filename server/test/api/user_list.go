package api

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ListUser() {
	testCase("operator not admin", listUserCase_NotAdmin)
	testCase("success", listUserCase_Success)
}

func listUserParams() string {
	return `{"page":{"size":10,"num":1}}`
}

func listUserCase_NotAdmin() string {
	res := httpInvoke(api.URI_ListUser, listUserParams(), accessToken_User)
	if res.IsSuccess || !errorIs(res.Err, utils.ErrPermissionDenied()) {
		return unknownError
	}

	return ""
}

func listUserCase_Success() string {
	res := httpInvoke(api.URI_ListUser, listUserParams(), accessToken_Admin)
	if !res.IsSuccess {
		return res.Err
	}

	count, _, err := dal.ListUsers(10, 1)
	if count != 4 || err != nil {
		return unknownError
	}

	return ""
}
