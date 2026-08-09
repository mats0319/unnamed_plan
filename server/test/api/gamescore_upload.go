package api

import (
	"fmt"

	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func UploadGameScore() {
	testCase("invalid game name", uploadGameScoreCase_InvalidGameName)
	testCase("success - visitor", uploadGameScoreCase_SuccessVisitor)
	testCase("success - user", uploadGameScoreCase_SuccessUser)
}

func uploadGameScoreParams(gameName api.GameName, score int, result string, player string) string {
	return fmt.Sprintf(`{"game_name":%d,"score":%d,"result":"%s","player":"%s"}`,
		gameName, score, result, player)
}

func uploadGameScoreCase_InvalidGameName() string {
	res := httpInvoke(api.URI_UploadGameScore, uploadGameScoreParams(-1, 100000, "test game result", "Visotor001"), "")
	if res.IsSuccess || !errorIs(res.Err, utils.ErrInvalidGameName()) {
		return unknownError
	}

	return ""
}

func uploadGameScoreCase_SuccessVisitor() string {
	res := httpInvoke(api.URI_UploadGameScore, uploadGameScoreParams(api.GameName_Flip, 100000, "test game result", "Visotor001"), "")
	if !res.IsSuccess {
		return res.Err
	}

	count, _, err := dal.ListFlipGameScores(10, 1)
	if err != nil || count != 1 {
		return unknownError
	}

	return ""
}

func uploadGameScoreCase_SuccessUser() string {
	res := httpInvoke(api.URI_UploadGameScore, uploadGameScoreParams(api.GameName_Flip, 100000, "test game result", ""), accessToken_User)
	if !res.IsSuccess {
		return res.Err
	}

	count, _, err := dal.ListFlipGameScores(10, 1)
	if err != nil || count != 2 {
		return unknownError
	}

	return ""
}
