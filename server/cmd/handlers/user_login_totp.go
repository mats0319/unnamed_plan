package handlers

import (
	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/cmd/config"
	"github.com/mats0319/unnamed_plan/server/cmd/handlers/mfa"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	"github.com/mats0319/unnamed_plan/server/internal/http/middleware"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func LoginMFA(ctx *mhttp.Context) {
	req := &api.LoginMFAReq{}
	if !ctx.ParseParams(req) {
		return
	}

	if len(req.MFAToken) < 1 || len(req.TOTPCode) < 1 {
		e := utils.ErrInvalidParams().WithParam("mfa token", req.MFAToken).WithParam("totp code", req.TOTPCode)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	t, e := mfa.VerifyMFAToken(req.MFAToken)
	if e != nil {
		ctx.ResData = e
		return
	}

	user, e := dal.GetUser(t.UserName)
	if e != nil {
		ctx.ResData = e
		return
	}

	decryptedKey, e := utils.Decrypt(user.TOTPKey, config.GetConfig().EncryptKey)
	if e != nil {
		ctx.ResData = e
		return
	}

	e = mfa.VerifyTOTPCode(req.TOTPCode, decryptedKey)
	if e != nil {
		ctx.ResData = e
		return
	}

	_ = dal.UpdateUser(user) // update user.UpdatedAt

	token := middleware.GenerateAPIAccessToken(user.UserName, config.GetConfig().AccessTokenExpireHour)
	ctx.SetHeader(utils.HTTPHeader_AccessToken, token)

	ctx.ResData = &api.LoginMFARes{
		UserName: user.UserName,
		Nickname: user.Nickname,
		IsAdmin:  user.IsAdmin,
	}
}
