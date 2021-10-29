package handlers

import (
	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/unnamed_plan/admin_data/http/structure_defination"
	"github.com/mats9693/utils/toy_server/http"
	"net/http"
)

func ListThinkingNoteByWriter(r *http.Request) *mhttp.ResponseData {
	params := &structure.ListThinkingNoteByWriterReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	if len(params.OperatorID) < 1 || params.PageSize < 1 || params.PageNum < 1 {
		return mhttp.ResponseWithError(error_InvalidParams,
			String("operator id", params.OperatorID),
			Int("page size", params.PageSize),
			Int("page num", params.PageNum))
	}

	notes, count, err := dao.GetThinkingNote().QueryPageByWriter(params.PageSize, params.PageNum, params.OperatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	noteListRes := make([]*structure.NoteRes, 0, len(notes))
	for i := range notes {
		noteListRes = append(noteListRes, noteDBToHTTPRes(notes[i]))
	}

	return mhttp.Response(structure.MakeListThinkingNoteByWriterRes(count, noteListRes))
}

func ListPublicThinkingNote(r *http.Request) *mhttp.ResponseData {
	params := &structure.ListPublicThinkingNoteReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	if len(params.OperatorID) < 1 || params.PageSize < 1 || params.PageNum < 1 {
		return mhttp.ResponseWithError(error_InvalidParams,
			String("operator id", params.OperatorID),
			Int("page size", params.PageSize),
			Int("page num", params.PageNum))
	}

	notes, count, err := dao.GetThinkingNote().QueryPageInPublic(params.PageSize, params.PageNum, params.OperatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	noteListRes := make([]*structure.NoteRes, 0, len(notes))
	for i := range notes {
		noteListRes = append(noteListRes, noteDBToHTTPRes(notes[i]))
	}

	return mhttp.Response(structure.MakeListPublicThinkingNoteRes(count, noteListRes))
}

func CreateThinkingNote(r *http.Request) *mhttp.ResponseData {
	params := &structure.CreateThinkingNoteReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	if len(params.OperatorID) < 1 || len(params.Content) < 1 {
		return mhttp.ResponseWithError(error_InvalidParams,
			String("operator id", params.OperatorID),
			String("topic(allow null)", params.Topic),
			String("content", params.Content))
	}

	err := dao.GetThinkingNote().Insert(&model.ThinkingNote{
		WriteBy:  params.OperatorID,
		Topic:    params.Topic,
		Content:  params.Content,
		IsPublic: params.IsPublic,
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeCreateThinkingNoteRes(true))
}

func ModifyThinkingNote(r *http.Request) *mhttp.ResponseData {
	params := &structure.ModifyThinkingNoteReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	if len(params.OperatorID) < 1 || len(params.NoteID) < 1 {
		return mhttp.ResponseWithError(error_InvalidParams,
			String("operator id", params.OperatorID),
			String("note id", params.NoteID))
	}

	_, err := verifyPwdByUserID(params.Password, params.OperatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	noteRecord, err := dao.GetThinkingNote().QueryFirst(model.ThinkingNote_NoteID+" = ?", params.NoteID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}
	if noteRecord.WriteBy != params.OperatorID {
		return mhttp.ResponseWithError(error_ModifyOthersThinkingNote)
	}

	if len(params.Topic)+len(params.Content) < 1 && noteRecord.IsPublic == params.IsPublic {
		return mhttp.ResponseWithError(error_NoValidModification)
	}

	updateColumns := make([]string, 0, 3)
	if len(params.Topic) > 0 {
		noteRecord.Topic = params.Topic
		updateColumns = append(updateColumns, model.ThinkingNote_Topic)
	}
	if len(params.Content) > 0 {
		noteRecord.Content = params.Content
		updateColumns = append(updateColumns, model.ThinkingNote_Content)
	}
	if noteRecord.IsPublic != params.IsPublic {
		noteRecord.IsPublic = params.IsPublic
		updateColumns = append(updateColumns, model.ThinkingNote_IsPublic)
	}

	err = dao.GetThinkingNote().UpdateColumnsByNoteID(noteRecord, updateColumns...)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeModifyThinkingNoteRes(true))
}

func DeleteThinkingNote(r *http.Request) *mhttp.ResponseData {
	params := &structure.DeleteThinkingNoteReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	if len(params.OperatorID) < 1 || len(params.Password) < 1 || len(params.NoteID) < 1 {
		return mhttp.ResponseWithError(error_InvalidParams,
			String("operator id", params.OperatorID),
			String("password", params.Password),
			String("note id", params.NoteID))
	}

	_, err := verifyPwdByUserID(params.Password, params.OperatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	err = dao.GetThinkingNote().UpdateColumnsByNoteID(&model.ThinkingNote{
		NoteID:    params.NoteID,
		IsDeleted: true,
	}, model.ThinkingNote_IsDeleted)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeDeleteThinkingNoteRes(true))
}

func noteDBToHTTPRes(data *model.ThinkingNote) *structure.NoteRes {
	if data == nil {
		return &structure.NoteRes{}
	}

	return &structure.NoteRes{
		NoteID:      data.NoteID,
		WriteBy:     data.WriteBy,
		Topic:       data.Topic,
		Content:     data.Content,
		IsPublic:    data.IsPublic,
		UpdateTime:  data.UpdateTime,
		CreatedTime: data.CreatedTime,
	}
}
