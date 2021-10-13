package handlers

import (
	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/unnamed_plan/admin_data/kits"
	"github.com/pkg/errors"
)

const (
	error_InvalidAccountOrPassword = "invalid account or password"
	error_InvalidParams            = "invalid param(s)"
	error_PermissionDenied         = "permission denied"
	error_NoValidModification      = "not any valid modification received"

	error_UserLocked   = "user is already locked"
	error_UserUnlocked = "user is already unlocked"

	error_ModifyOthersThinkingNote = "not allowed to modify others' thinking note"
)

func checkPwdByUserName(password string, userName string) (user *model.User, err error) {
	user, err = dao.GetUser().QueryOneInUnlocked(model.User_UserName+" = ?", userName)
	if err != nil {
		return
	}

	if user.Password != kits.CalcSHA256(password, user.Salt) {
		err = errors.New(error_InvalidAccountOrPassword)
		return
	}

	return
}

func checkPwdByUserID(password string, userID string) (user *model.User, err error) {
	user, err = dao.GetUser().QueryOneInUnlocked(model.User_UserID+" = ?", userID)
	if err != nil {
		return
	}

	if user.Password != kits.CalcSHA256(password, user.Salt) {
		err = errors.New(error_InvalidAccountOrPassword)
		return
	}

	return
}

func errorsToString(errs ...error) string {
	res := ""
	for i := range errs {
		if errs[i] != nil {
			res += errs[i].Error()
		}
	}

	return res
}
