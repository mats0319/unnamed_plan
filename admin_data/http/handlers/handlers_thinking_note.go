package handlers

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/unnamed_plan/admin_data/kits"
	mhttp "github.com/mats9693/utils/toy_server/http"
	"net/http"
)

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

func ListThinkingNoteByWriter(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
}

func ListPublicThinkingNote(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
}

func ListDeletedThinkingNote(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
}

func DeleteThinkingNote(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
}
