package testdata

import "github.com/mats0319/unnamed_plan/server/internal/db/model"

var TestFlipGameScores = []*model.FlipGameScore{
	{
		GameScore: model.GameScore{
			Score:      10000,
			Result:     `{"duration":20,"steps":18}`,
			Player:     "mats0319",
			PlayerName: "Mario",
		},
	},
	{
		GameScore: model.GameScore{
			Score:      9000,
			Result:     `{"duration":30,"steps":25}`,
			Player:     "mats0319",
			PlayerName: "Mario",
		},
	},
	{
		GameScore: model.GameScore{
			Score:      8000,
			Result:     `{"duration":40,"steps":30}`,
			Player:     "",
			PlayerName: "visitor 001",
		},
	},
}
