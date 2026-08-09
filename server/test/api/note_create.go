package api

import (
	"context"
	"fmt"

	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

var noteID string

func CreateNote() {
	testCase("success", createNoteCase_Success)
	testCase("duplicate create", createNoteCase_Duplicate)
}

func createNoteParams(isAnonymous bool, title string, content string) string {
	return fmt.Sprintf(`{"is_anonymous":%t,"title":"%s","content":"%s"}`, isAnonymous, title, content)
}

func createNoteCase_Success() string {
	res := httpInvoke(api.URI_CreateNote, createNoteParams(false, "123", "456"), accessToken_Admin)
	if !res.IsSuccess {
		return res.Err
	}

	// record 'note id' for later api(s)
	note, _ := dal.Note.WithContext(context.TODO()).First()
	if note == nil {
		return unknownError
	}
	noteID = note.NoteID

	data, err := dal.GetNote(noteID)
	if err != nil || data == nil || data.Title != "123" || data.Content != "456" {
		return unknownError
	}

	return ""
}

func createNoteCase_Duplicate() string {
	res := httpInvoke(api.URI_CreateNote, createNoteParams(false, "123", "456"), accessToken_Admin)
	if res.IsSuccess || !errorIs(res.Err, utils.ErrNoteExist()) {
		return unknownError
	}

	return ""
}
