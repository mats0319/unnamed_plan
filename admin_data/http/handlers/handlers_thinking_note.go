package handlers

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
    "github.com/mats9693/unnamed_plan/admin_data/http/structure_defination"
    "github.com/mats9693/unnamed_plan/admin_data/utils"
	"github.com/mats9693/utils/toy_server/http"
	"net/http"
	"strconv"
)

func ListThinkingNoteByWriter(r *http.Request) string {
	operatorID := r.PostFormValue("operatorID")
	pageSize, err := strconv.Atoi(r.PostFormValue("pageSize"))
	pageNum, err2 := strconv.Atoi(r.PostFormValue("pageNum"))

	if err != nil || err2 != nil {
		return mhttp.ResponseWithError(errorsToString(err, err2))
	}
	if len(operatorID) < 1 || pageSize < 1 || pageNum < 1 {
		return mhttp.ResponseWithError(error_InvalidParams+
			fmt.Sprintf(", operator id: %s, page size: %d, page num: %d", operatorID, pageSize, pageNum))
	}

	notes, count, err := dao.GetThinkingNote().QueryPageByWriter(pageSize, pageNum, operatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	noteListRes := make([]*structure.NoteListRes, 0, len(notes))
	for i := range notes {
		noteListRes = append(noteListRes, &structure.NoteListRes{
			NoteID:      notes[i].NoteID,
			WriteBy:     notes[i].WriteBy,
			Topic:       notes[i].Topic,
			Content:     notes[i].Content,
			IsPublic:    notes[i].IsPublic,
			UpdateTime:  notes[i].UpdateTime,
			CreatedTime: notes[i].CreatedTime,
		})
	}

	return mhttp.Response(structure.MakeListThinkingNoteByWriterRes(count, noteListRes))
}

func ListPublicThinkingNote(r *http.Request) string {
	operatorID := r.PostFormValue("operatorID")
	pageSize, err := strconv.Atoi(r.PostFormValue("pageSize"))
	pageNum, err2 := strconv.Atoi(r.PostFormValue("pageNum"))

	if err != nil || err2 != nil {
		return mhttp.ResponseWithError(errorsToString(err, err2))
	}
	if len(operatorID) < 1 || pageSize < 1 || pageNum < 1 {
		return mhttp.ResponseWithError(error_InvalidParams+
			fmt.Sprintf(", operator id: %s, page size: %d, page num: %d", operatorID, pageSize, pageNum))
	}

	notes, count, err := dao.GetThinkingNote().QueryPageInPublic(pageSize, pageNum, operatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	noteListRes := make([]*structure.NoteListRes, 0, len(notes))
	for i := range notes {
		noteListRes = append(noteListRes, &structure.NoteListRes{
			NoteID:      notes[i].NoteID,
			WriteBy:     notes[i].WriteBy,
			Topic:       notes[i].Topic,
			Content:     notes[i].Content,
			IsPublic:    notes[i].IsPublic,
			UpdateTime:  notes[i].UpdateTime,
			CreatedTime: notes[i].CreatedTime,
		})
	}

	return mhttp.Response(structure.MakeListPublicThinkingNoteRes(count, noteListRes))
}

func CreateThinkingNote(r *http.Request) string {
	operatorID := r.PostFormValue("operatorID")
	topic := r.PostFormValue("topic")
	content := r.PostFormValue("content")
	isPublic, err := utils.StringToBool(r.PostFormValue("isPublic"))

	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}
	if len(operatorID) < 1 || len(content) < 1 {
		return mhttp.ResponseWithError(error_InvalidParams+
			fmt.Sprintf(", operator id: %s, content length: %d", operatorID, len(content)))
	}

	err = dao.GetThinkingNote().Insert(&model.ThinkingNote{
		WriteBy:  operatorID,
		Topic:    topic,
		Content:  content,
		IsPublic: isPublic,
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeCreateThinkingNoteRes(true))
}

func ModifyThinkingNote(r *http.Request) string {
	operatorID := r.PostFormValue("operatorID")
	noteID := r.PostFormValue("noteID")
	password := r.PostFormValue("password")
	topic := r.PostFormValue("topic")
	content := r.PostFormValue("content")
	isPublicStr := r.PostFormValue("isPublic")

	if len(operatorID) < 1 || len(noteID) < 1 {
		return mhttp.ResponseWithError(error_InvalidParams+
			fmt.Sprintf(", operator id: %s, note id: %s", operatorID, noteID))
	}
	if len(topic)+len(content)+len(isPublicStr) < 1 {
		return mhttp.ResponseWithError(error_NoValidModification)
	}

	isPublic, err := utils.StringToBool(isPublicStr)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	_, err = verifyPwdByUserID(password, operatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	newNote, err := dao.GetThinkingNote().QueryFirst(model.ThinkingNote_NoteID+" = ?", noteID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}
	if newNote.WriteBy != operatorID {
		return mhttp.ResponseWithError(error_ModifyOthersThinkingNote)
	}

	updateColumns := make([]string, 0, 3)
	if len(topic) > 0 {
		newNote.Topic = topic
		updateColumns = append(updateColumns, model.ThinkingNote_Topic)
	}
	if len(content) > 0 {
		newNote.Content = content
		updateColumns = append(updateColumns, model.ThinkingNote_Content)
	}
	if len(isPublicStr) > 0 {
		newNote.IsPublic = isPublic
		updateColumns = append(updateColumns, model.ThinkingNote_IsPublic)
	}

	err = dao.GetThinkingNote().UpdateColumnsByNoteID(newNote, updateColumns...)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeModifyThinkingNoteRes(true))
}

func DeleteThinkingNote(r *http.Request) string {
	operatorID := r.PostFormValue("operatorID")
	password := r.PostFormValue("password")
	noteID := r.PostFormValue("noteID")

	if len(operatorID) < 1 || len(password) < 1 || len(noteID) < 1 {
		return mhttp.ResponseWithError(error_InvalidParams+
			fmt.Sprintf(", operator id: %s, password: %s, note id: %s", operatorID, password, noteID))
	}

	_, err := verifyPwdByUserID(password, operatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	err = dao.GetThinkingNote().UpdateColumnsByNoteID(&model.ThinkingNote{
		NoteID:    noteID,
		IsDeleted: true,
	}, model.ThinkingNote_IsDeleted)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeDeleteThinkingNoteRes(true))
}
