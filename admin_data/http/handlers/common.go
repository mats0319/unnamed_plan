package handlers

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/admin_data/utils"
	"github.com/mats9693/unnamed_plan/shared/db/dao"
	"github.com/mats9693/unnamed_plan/shared/db/model"
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

func verifyPwdByUserName(password string, userName string) (user *model.User, err error) {
	user, err = dao.GetUser().QueryOneInUnlocked(model.User_UserName+" = ?", userName)
	if err != nil {
		return
	}

	if user.Password != utils.CalcSHA256(password, user.Salt) {
		err = errors.New(error_InvalidAccountOrPassword)
		return
	}

	return
}

func verifyPwdByUserID(password string, userID string) (user *model.User, err error) {
	user, err = dao.GetUser().QueryOneInUnlocked(model.User_UserID+" = ?", userID)
	if err != nil {
		return
	}

	if user.Password != utils.CalcSHA256(password, user.Salt) {
		err = errors.New(error_InvalidAccountOrPassword)
		return
	}

	return
}

func String(key string, value string) string {
	return fmt.Sprintf(" { %s : %s } ", key, value)
}

func Int(key string, value int) string {
	return fmt.Sprintf(" { %s : %d } ", key, value)
}

func Uint8(key string, value uint8) string {
	return fmt.Sprintf(" { %s : %d } ", key, value)
}
