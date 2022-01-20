package handlers

import (
    "github.com/mats9693/unnamed_plan/services/config_center/db"
    "github.com/mats9693/unnamed_plan/services/config_center/http/structure_defination"
    "github.com/mats9693/unnamed_plan/services/shared/const"
    "github.com/mats9693/unnamed_plan/services/shared/http"
    "github.com/mats9693/unnamed_plan/services/shared/log"
    "github.com/mats9693/unnamed_plan/services/shared/utils"
    "go.uber.org/zap"
    "net/http"
)

func Login(r *http.Request) *mhttp.ResponseData {
    params := &structure.LoginReqParams{}
    if errMsg := params.Decode(r); len(errMsg) > 0 {
        mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
        return mhttp.ResponseWithError(errMsg)
    }

    if len(params.UserName) < 1 || len(params.Password) < 1 {
        return mhttp.ResponseWithError(mconst.Error_InvalidParams)
    }

    user, err := db.GetAuthenticateDao().QueryOneByUserName(params.UserName)
    if err != nil {
        return mhttp.ResponseWithError(err.Error())
    }

    if user.Password != utils.CalcSHA256(params.Password, user.Salt) {
        return mhttp.ResponseWithError(mconst.Error_InvalidAccountOrPassword)
    }

    return mhttp.Response(structure.MakeLoginRes(user.ID, user.UserName), user.ID)
}
