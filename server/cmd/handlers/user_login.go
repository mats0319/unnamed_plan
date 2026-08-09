package handlers

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/cmd/config"
	"github.com/mats0319/unnamed_plan/server/cmd/handlers/mfa"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	"github.com/mats0319/unnamed_plan/server/internal/http/middleware"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/password"
)

func Login(ctx *mhttp.Context) {
	req := &api.LoginReq{}
	if !ctx.ParseParams(req) {
		return
	}

	if len(req.UserName) < 1 || len(req.Password) < 1 {
		e := utils.ErrInvalidParams().WithParam("user name", req.UserName).WithParam("password", req.Password)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	user, e := dal.GetUser(req.UserName)
	if e != nil {
		ctx.ResData = e
		return
	}

	e = password.VerifyPassword(req.Password, user.Password)
	if e != nil {
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	// 如果启用MFA：暂停登录，进入MFA过程
	if user.EnableMFA {
		ctx.ResData = &api.LoginRes{
			EnableMFA: true,
			MFAToken:  mfa.GenerateMFAToken(user.UserName, config.GetConfig().MFATokenExpireMinute),
		}

		return
	}

	_ = dal.UpdateUser(user) // update user.UpdatedAt

	token := middleware.GenerateAPIAccessToken(user.UserName, config.GetConfig().AccessTokenExpireHour)
	ctx.SetHeader(utils.HTTPHeader_AccessToken, token)

	ctx.ResData = &api.LoginRes{
		EnableMFA:  false,
		UserName:   user.UserName,
		Nickname:   user.Nickname,
		IsAdmin:    user.IsAdmin,
		HasTOTPKey: len(user.TOTPKey) > 0,
	}
}
