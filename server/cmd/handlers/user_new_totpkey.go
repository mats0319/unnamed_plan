package handlers

import (
	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/cmd/config"
	"github.com/mats0319/unnamed_plan/server/cmd/handlers/mfa"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func NewTOTPKey(ctx *mhttp.Context) {
	if len(ctx.UserName) < 1 {
		e := utils.ErrInvalidParams().WithParam("operator", ctx.UserName)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	totpKey, e := mfa.GenerateTOTPKey(
		ctx.UserName,
		config.GetConfig().TOTPKeyExpireMinute,
		config.GetConfig().EncryptKey,
	)
	if e != nil {
		ctx.ResData = e
		return
	}

	ctx.ResData = &api.NewTOTPKeyRes{TOTPKey: totpKey}
}
