package handlers

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ListUser(ctx *mhttp.Context) {
	req := &api.ListUserReq{}
	if !ctx.ParseParams(req) {
		return
	}

	if len(ctx.UserName) < 1 || req.Page.Size <= 0 || req.Page.Num <= 0 {
		e := utils.ErrInvalidParams().WithParam("operator", ctx.UserName).WithParam("pagination", req.Page)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	operator, e := dal.GetUser(ctx.UserName)
	if e != nil {
		ctx.ResData = e
		return
	}

	if !operator.IsAdmin {
		e := utils.ErrPermissionDenied().WithParam("need admin", operator.UserName)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	count, users, e := dal.ListUsers(req.Page.Size, req.Page.Num)
	if e != nil {
		ctx.ResData = e
		return
	}

	ctx.ResData = &api.ListUserRes{
		Count: count,
		Users: usersDBToHTTP(users),
	}
}

func usersDBToHTTP(users []*model.User) []*api.User {
	res := make([]*api.User, len(users))
	for i, v := range users {
		res[i] = &api.User{
			UserName:   v.UserName,
			Nickname:   v.Nickname,
			CreatedAt:  v.CreatedAt,
			IsAdmin:    v.IsAdmin,
			EnableMFA:  v.EnableMFA,
			HasTOTPKey: len(v.TOTPKey) > 0,
			LastLogin:  v.UpdatedAt,
		}
	}

	return res
}
