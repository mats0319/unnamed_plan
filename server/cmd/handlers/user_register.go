package handlers

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/password"
)

func Register(ctx *mhttp.Context) {
	req := &api.RegisterReq{}
	if !ctx.ParseParams(req) {
		return
	}

	if len(req.UserName) < 1 || len(req.Password) < 1 {
		e := utils.ErrInvalidParams().WithParam("user name", req.UserName).WithParam("password", req.Password)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	user := &model.User{
		UserName: req.UserName,
		Nickname: req.UserName,
		Password: password.GeneratePassword(req.Password),
	}

	e := dal.CreateUser(user)
	if e != nil {
		ctx.ResData = e
		return
	}

	ctx.ResData = &api.RegisterRes{}
}
