package dal

import (
	"context"

	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func CreateFlipGameScore(gameScore *model.FlipGameScore) *utils.Error {
	err := FlipGameScore.WithContext(context.TODO()).Create(gameScore)
	if err != nil {
		e := utils.ErrDBError().WithCause(err) // 游戏成绩表的字段均没有唯一约束，所以这里不需要把唯一约束冲突错误单独拎出来
		mlog.Error(e.String())
		return e
	}

	if len(gameScore.Player) > 0 { // 玩家已登录
		records, err := FlipGameScore.WithContext(context.TODO()).Where(FlipGameScore.Player.Eq(gameScore.Player)).
			Order(FlipGameScore.Score.Desc()).Order(FlipGameScore.CreatedAt.Asc()).Find()
		if err != nil {
			e := utils.ErrDBError().WithCause(err)
			mlog.Error(e.String())
			return e
		}
		if len(records) > 3 { // 如果一个玩家拥有超过3条成绩，只保留分数最高的3条
			_, err := FlipGameScore.WithContext(context.TODO()).Delete(records[3:]...)
			if err != nil {
				e := utils.ErrDBError().WithCause(err)
				mlog.Error(e.String())
				return e
			}
		}
	}

	return nil
}

func ListFlipGameScores(pageSize int, pageNum int) (count int64, records []*model.FlipGameScore, e *utils.Error) {
	sql := FlipGameScore.WithContext(context.TODO())

	count, err := sql.Count()
	if err != nil {
		e = utils.ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return
	}

	records, err = sql.Order(FlipGameScore.Score.Desc()).Offset(pageSize * (pageNum - 1)).Limit(pageSize).Find()
	if err != nil {
		e = utils.ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return
	}

	return
}
