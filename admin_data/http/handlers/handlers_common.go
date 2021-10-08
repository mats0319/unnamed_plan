package handlers

import (
    "github.com/mats9693/unnamed_plan/admin_data/db/dao"
    "github.com/mats9693/unnamed_plan/admin_data/db/model"
    "github.com/mats9693/unnamed_plan/admin_data/kits"
    "github.com/pkg/errors"
)

func checkPwdByUserName(password string, userName string) (user *model.User, err error) {
    user, err = dao.GetUser().QueryUnlocked(model.User_UserName+" = ?", userName)
    if err != nil {
        return
    }

    if user.Password != kits.CalcSHA256(password, user.Salt) {
        err = errors.New("invalid account or password")
        return
    }

    return
}

func checkPwdByUserID(password string, userID string) (user *model.User, err error) {
    user, err = dao.GetUser().QueryUnlocked(model.User_UserID+" = ?", userID)
    if err != nil {
        return
    }

    if user.Password != kits.CalcSHA256(password, user.Salt) {
        err = errors.New("invalid account or password")
        return
    }

    return
}
