package api

import (
	"fmt"

	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
)

func ListNote() {
	testCase("success - list all", listNoteCase_SuccessListAll)
	testCase("success - list only operator", listNoteCase_SuccessListOnlyOperator)
}

func listNoteParams(onlyOperator bool) string {
	return fmt.Sprintf(`{"page":{"size":10,"num":1},"only_operator":%t}`, onlyOperator)
}

func listNoteCase_SuccessListAll() string {
	res := httpInvoke(api.URI_ListNote, listNoteParams(false), accessToken_User)
	if !res.IsSuccess {
		return res.Err
	}

	count, _, err := dal.ListNotes(10, 1, "")
	if err != nil || count != 1 {
		return unknownError
	}

	return ""
}

func listNoteCase_SuccessListOnlyOperator() string {
	res := httpInvoke(api.URI_ListNote, listNoteParams(true), accessToken_User)
	if !res.IsSuccess {
		return res.Err
	}

	count, _, err := dal.ListNotes(10, 1, "not exist")
	if err != nil || count != 0 {
		return unknownError
	}

	return ""
}
