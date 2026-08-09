package handlers

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/password"
)

func ModifyUser(ctx *mhttp.Context) {
	req := &api.ModifyUserReq{}
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

	modifyNicknameFlag := len(req.Nickname) > 0 && req.Nickname != operator.Nickname
	if !modifyNicknameFlag && len(req.Password) < 1 {
		e := utils.ErrNoChanges().WithParam("operator", operator.UserName)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	// modify
	if modifyNicknameFlag {
		operator.Nickname = req.Nickname
	}
	if len(req.Password) > 0 {
		// when modify, same pwd is wrong
		if password.VerifyPassword(req.Password, operator.Password) == nil {
			e := utils.ErrSamePassword()
			mlog.Error(e.String())
			ctx.ResData = e
			return
		}

		operator.Password = password.GeneratePassword(req.Password)
	}

	e = dal.UpdateUser(operator)
	if e != nil {
		ctx.ResData = e
		return
	}

	ctx.ResData = &api.ModifyUserRes{}
}
