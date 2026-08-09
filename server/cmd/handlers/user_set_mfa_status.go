package handlers

import (
	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/cmd/config"
	"github.com/mats0319/unnamed_plan/server/cmd/handlers/mfa"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func SetMFAStatus(ctx *mhttp.Context) {
	req := &api.SetMFAStatusReq{}
	if !ctx.ParseParams(req) {
		return
	}

	if len(ctx.UserName) < 1 {
		e := utils.ErrInvalidParams().WithParam("operator", ctx.UserName)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	operator, e := dal.GetUser(ctx.UserName)
	if e != nil {
		ctx.ResData = e
		return
	}

	/* 禁用MFA */

	if !req.EnableMFA {
		if operator.EnableMFA { // 仅在启用状态下更新数据库，如已是禁用状态，直接返回
			operator.EnableMFA = false
			e := dal.UpdateUser(operator)
			if e != nil {
				ctx.ResData = e
				return
			}
		}

		ctx.ResData = &api.SetMFAStatusRes{}

		return
	}

	/* 启用MFA */

	if req.ApplyNewKeyFlag { // 申请了新的totp key
		encryptedKeyBytes, e := mfa.VerifyTOTPKey(ctx.UserName, req.TOTPCode, config.GetConfig().EncryptKey)
		if e != nil {
			ctx.ResData = e
			return
		}
		operator.TOTPKey = encryptedKeyBytes
	} else { // 使用历史totp key
		decryptedKey, e := utils.Decrypt(operator.TOTPKey, config.GetConfig().EncryptKey)
		if e != nil {
			ctx.ResData = e
			return
		}

		e = mfa.VerifyTOTPCode(req.TOTPCode, decryptedKey)
		if e != nil {
			ctx.ResData = e
			return
		}
	}

	operator.EnableMFA = true
	e = dal.UpdateUser(operator)
	if e != nil {
		ctx.ResData = e
		return
	}

	ctx.ResData = &api.SetMFAStatusRes{}
}
