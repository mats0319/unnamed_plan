package handlers

import (
	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ListGameScore(ctx *mhttp.Context) {
	req := &api.ListGameScoreReq{}
	if !ctx.ParseParams(req) {
		return
	}

	if req.Page.Size <= 0 || req.Page.Num <= 0 {
		e := utils.ErrInvalidParams().WithParam("pagination", req.Page)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	var (
		count         int64
		gameScoreHTTP []*api.GameScore
		e             *utils.Error
	)
	switch req.GameName {
	case api.GameName_Flip:
		records := make([]*model.FlipGameScore, 0)
		count, records, e = dal.ListFlipGameScores(req.Page.Size, req.Page.Num)
		if e != nil {
			ctx.ResData = e
			return
		}
		gameScoreHTTP = flipGameScoresDBToHTTP(records)
	default:
		e = utils.ErrInvalidGameName().WithParam("game name", req.GameName)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	ctx.ResData = &api.ListGameScoreRes{Count: count, Scores: gameScoreHTTP}
}

func flipGameScoresDBToHTTP(scores []*model.FlipGameScore) []*api.GameScore {
	res := make([]*api.GameScore, len(scores))
	for i, v := range scores {
		res[i] = &api.GameScore{
			Score:      v.Score,
			Result:     v.Result,
			PlayerName: v.PlayerName,
			Timestamp:  v.CreatedAt,
		}
	}

	return res
}
