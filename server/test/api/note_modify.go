package api

import (
	"fmt"

	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ModifyNote() {
	testCase("note not exist", modifyNoteCase_NoteNotExist)
	testCase("no changes", modifyNoteCase_NoChanges)
	testCase("not writer", modifyNoteCase_NotWriter)
	testCase("success", modifyNoteCase_Success)
}

func modifyNoteParams(noteID string, isAnonymous bool, title string, content string) string {
	return fmt.Sprintf(`{"note_id":"%s","is_anonymous":%t,"title":"%s","content":"%s"}`,
		noteID, isAnonymous, title, content)
}

func modifyNoteCase_NoteNotExist() string {
	res := httpInvoke(api.URI_ModifyNote, modifyNoteParams("not exist", false, "", "1"), accessToken_User)
	if res.IsSuccess || !errorIs(res.Err, utils.ErrNoteNotFound()) {
		return unknownError
	}

	return ""
}

func modifyNoteCase_NoChanges() string {
	res := httpInvoke(api.URI_ModifyNote, modifyNoteParams(noteID, false, "123", "456"), accessToken_Admin)
	if res.IsSuccess || !errorIs(res.Err, utils.ErrNoChanges()) {
		return unknownError
	}

	return ""
}

func modifyNoteCase_NotWriter() string {
	res := httpInvoke(api.URI_ModifyNote, modifyNoteParams(noteID, false, "", "1"), accessToken_User)
	if res.IsSuccess || !errorIs(res.Err, utils.ErrPermissionDenied()) {
		return unknownError
	}

	return ""
}

func modifyNoteCase_Success() string {
	res := httpInvoke(api.URI_ModifyNote, modifyNoteParams(noteID, true, "123123", "456456"), accessToken_Admin)
	if !res.IsSuccess {
		return res.Err
	}

	data, err := dal.GetNote(noteID)
	if data == nil || !data.IsAnonymous || data.Title != "123123" || data.Content != "456456" || err != nil {
		return unknownError
	}

	return ""
}
