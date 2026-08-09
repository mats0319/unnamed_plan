package api

import (
	"fmt"

	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func DeleteNote() {
	testCase("note not exist", deleteNoteCase_NoteNotExist)
	testCase("not writer", deleteNoteCase_NotWriter)
	testCase("success", deleteNoteCase_Success)
}

func deleteNoteParams(noteID string) string {
	return fmt.Sprintf(`{"note_id":"%s"}`, noteID)
}

func deleteNoteCase_NoteNotExist() string {
	res := httpInvoke(api.URI_DeleteNote, deleteNoteParams("not exist"), accessToken_User)
	if res.IsSuccess || !errorIs(res.Err, utils.ErrNoteNotFound()) {
		return unknownError
	}

	return ""
}

func deleteNoteCase_NotWriter() string {
	res := httpInvoke(api.URI_DeleteNote, deleteNoteParams(noteID), accessToken_User)
	if res.IsSuccess || !errorIs(res.Err, utils.ErrPermissionDenied()) {
		return unknownError
	}

	return ""
}

func deleteNoteCase_Success() string {
	res := httpInvoke(api.URI_DeleteNote, deleteNoteParams(noteID), accessToken_Admin)
	if !res.IsSuccess {
		return res.Err
	}

	count, _, err := dal.ListNotes(10, 1, "")
	if err != nil || count != 0 {
		return unknownError
	}

	return ""
}
