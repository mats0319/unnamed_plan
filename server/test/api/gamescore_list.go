package api

import (
	"fmt"

	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ListGameScore() {
	testCase("invalid game name", listGameScoreCase_InvalidGameName)
	testCase("success", listGameScoreCase_Success)
}

func listGameScoreParams(gameName api.GameName) string {
	return fmt.Sprintf(`{"game_name":%d,"page":{"size":10,"num":1}}`, gameName)
}

func listGameScoreCase_InvalidGameName() string {
	res := httpInvoke(api.URI_ListGameScore, listGameScoreParams(-1), "")
	if res.IsSuccess || !errorIs(res.Err, utils.ErrInvalidGameName()) {
		return unknownError
	}

	return ""
}

func listGameScoreCase_Success() string {
	res := httpInvoke(api.URI_ListGameScore, listGameScoreParams(api.GameName_Flip), "")
	if !res.IsSuccess {
		return res.Err
	}

	return ""
}
