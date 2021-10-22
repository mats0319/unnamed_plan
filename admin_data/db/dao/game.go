package dao

import "github.com/mats9693/unnamed_plan/admin_data/db/model"

type GameResult model.GameResult

var gameResultIns = &GameResult{}

func GetGameResult() *GameResult {
    return gameResultIns
}
