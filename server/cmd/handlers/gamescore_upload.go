package handlers

import (
	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func UploadGameScore(ctx *mhttp.Context) {
	req := &api.UploadGameScoreReq{}
	if !ctx.ParseParams(req) {
		return
	}

	gameScore := model.GameScore{Score: req.Score, Result: req.Result}

	// 填写获得该成绩的玩家信息
	if len(ctx.UserName) > 0 { // 已登录：填写用户名和昵称
		user, e := dal.GetUser(ctx.UserName)
		if e != nil {
			ctx.ResData = e
			return
		}
		gameScore.Player = user.UserName
		gameScore.PlayerName = user.Nickname
	} else { // 未登录：填写前端传的访客id
		gameScore.PlayerName = req.Player
	}

	var e *utils.Error
	switch req.GameName {
	case api.GameName_Flip:
		e = dal.CreateFlipGameScore(&model.FlipGameScore{GameScore: gameScore})
	default:
		e = utils.ErrInvalidGameName().WithParam("game name", req.GameName)
		mlog.Error(e.String())
	}
	if e != nil {
		ctx.ResData = e
		return
	}

	ctx.ResData = &api.UploadGameScoreRes{}
}
