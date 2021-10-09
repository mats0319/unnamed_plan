package handlers

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/unnamed_plan/admin_data/http/response_type"
	"github.com/mats9693/unnamed_plan/admin_data/kits"
	"github.com/mats9693/utils/toy_server/http"
	"net/http"
	"strconv"
)

func ListThinkingNoteByWriter(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	pageSize, err := strconv.Atoi(r.PostFormValue("pageSize"))
	pageNum, err2 := strconv.Atoi(r.PostFormValue("pageNum"))

	if err != nil || err2 != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(errorsToString(err, err2)))
		return
	}
	if len(operatorID) < 1 || pageSize < 1 || pageNum < 1 {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, page size: %d, page num: %d", operatorID, pageSize, pageNum)))
		return
	}

	notes, count, err := dao.GetThinkingNote().QueryPageByWriter(pageSize, pageNum, operatorID)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	notesRes := make([]*http_res_type.HTTPResNote, 0, len(notes))
	for i := range notes {
		notesRes = append(notesRes, &http_res_type.HTTPResNote{
			NoteID:      notes[i].NoteID,
			WriteBy:     notes[i].WriteBy,
			Topic:       notes[i].Topic,
			Content:     notes[i].Content,
			IsPublic:    notes[i].IsPublic,
			UpdateTime:  notes[i].UpdateTime,
			CreatedTime: notes[i].CreatedTime,
		})
	}

	resData := &struct {
		Total int                          `json:"total"`
		Notes []*http_res_type.HTTPResNote `json:"notes"`
	}{
		Total: count,
		Notes: notesRes,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}

func ListPublicThinkingNote(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	pageSize, err := strconv.Atoi(r.PostFormValue("pageSize"))
	pageNum, err2 := strconv.Atoi(r.PostFormValue("pageNum"))

	if err != nil || err2 != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(errorsToString(err, err2)))
		return
	}
	if len(operatorID) < 1 || pageSize < 1 || pageNum < 1 {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, page size: %d, page num: %d", operatorID, pageSize, pageNum)))
		return
	}

	notes, count, err := dao.GetThinkingNote().QueryPageInPublic(pageSize, pageNum, operatorID)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	notesRes := make([]*http_res_type.HTTPResNote, 0, len(notes))
	for i := range notes {
		notesRes = append(notesRes, &http_res_type.HTTPResNote{
			NoteID:      notes[i].NoteID,
			WriteBy:     notes[i].WriteBy,
			Topic:       notes[i].Topic,
			Content:     notes[i].Content,
			IsPublic:    notes[i].IsPublic,
			UpdateTime:  notes[i].UpdateTime,
			CreatedTime: notes[i].CreatedTime,
		})
	}

	resData := &struct {
		Total int                          `json:"total"`
		Notes []*http_res_type.HTTPResNote `json:"notes"`
	}{
		Total: count,
		Notes: notesRes,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}

func CreateThinkingNote(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	topic := r.PostFormValue("topic")
	content := r.PostFormValue("content")
	isPublic, err := kits.StringToBool(r.PostFormValue("isPublic"))

	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}
	if len(operatorID) < 1 || len(content) < 1 {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, content length: %d", operatorID, len(content))))
		return
	}

	err = dao.GetThinkingNote().Insert(&model.ThinkingNote{
		WriteBy:  operatorID,
		Topic:    topic,
		Content:  content,
		IsPublic: isPublic,
	})
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	resData := &struct {
		IsSuccess bool `json:"isSuccess"`
	}{
		IsSuccess: true,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}

func ModifyThinkingNote(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	noteID := r.PostFormValue("noteID")
	password := r.PostFormValue("password")
	topic := r.PostFormValue("topic")
	content := r.PostFormValue("content")
	isPublicStr := r.PostFormValue("isPublic")

	if len(operatorID) < 1 || len(noteID) < 1 {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, note id: %s", operatorID, noteID)))
		return
	}
	if len(topic)+len(content)+len(isPublicStr) < 1 {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError("invalid params, not any modification received"))
		return
	}

	isPublic, err := kits.StringToBool(isPublicStr)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	_, err = checkPwdByUserID(password, operatorID)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	newNote, err := dao.GetThinkingNote().QueryFirst(model.ThinkingNote_NoteID + " = ?", noteID)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}
	if newNote.WriteBy != operatorID {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError("invalid params, operator is not the writer of the note"))
		return
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
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	resData := &struct {
		IsSuccess bool `json:"isSuccess"`
	}{
		IsSuccess: true,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}

func DeleteThinkingNote(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	password := r.PostFormValue("password")
	noteID := r.PostFormValue("noteID")

	if len(operatorID) < 1 || len(password) < 1 || len(noteID) < 1 {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, password: %s, note id: %s", operatorID, password, noteID)))
		return
	}

	_, err := checkPwdByUserID(password, operatorID)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	err = dao.GetThinkingNote().UpdateColumnsByNoteID(&model.ThinkingNote{
		NoteID:    noteID,
		IsDeleted: true,
	}, model.ThinkingNote_IsDeleted)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	resData := &struct {
		IsSuccess bool `json:"isSuccess"`
	}{
		IsSuccess: true,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}
