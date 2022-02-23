package handlers

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/gateway/http/structure_defination"
	"github.com/mats9693/unnamed_plan/services/gateway/rpc"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/http"
	"github.com/mats9693/unnamed_plan/services/shared/http/response"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func ListNoteByWriter(r *http.Request) *mresponse.ResponseData {
	params := &structure.ListNoteByWriterReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	res, err := rpc.GetRPCClient().NoteClient.ListByWriter(context.Background(), &rpc_impl.Note_ListByWriterReq{
		OperatorId: params.OperatorID,
		Page: &rpc_impl.Pagination{
			PageSize: uint32(params.PageSize),
			PageNum:  uint32(params.PageNum),
		},
	})
	if err != nil {
		mlog.Logger().Error("list note by writer failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeListNoteByWriterRes(res.Total, notesRPCToHTTP(res.Notes...)))
}

func ListPublicNote(r *http.Request) *mresponse.ResponseData {
	params := &structure.ListPublicNoteReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	res, err := rpc.GetRPCClient().NoteClient.ListPublic(context.Background(), &rpc_impl.Note_ListPublicReq{
		OperatorId: params.OperatorID,
		Page: &rpc_impl.Pagination{
			PageSize: uint32(params.PageSize),
			PageNum:  uint32(params.PageNum),
		},
	})
	if err != nil {
		mlog.Logger().Error("list public note failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeListPublicNoteRes(res.Total, notesRPCToHTTP(res.Notes...)))
}

func CreateNote(r *http.Request) *mresponse.ResponseData {
	params := &structure.CreateNoteReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	_, err := rpc.GetRPCClient().NoteClient.Create(context.Background(), &rpc_impl.Note_CreateReq{
		OperatorId: params.OperatorID,
		Topic:      params.Topic,
		Content:    params.Content,
		IsPublic:   params.IsPublic,
	})
	if err != nil {
		mlog.Logger().Error("create note failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func ModifyNote(r *http.Request) *mresponse.ResponseData {
	params := &structure.ModifyNoteReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	_, err := rpc.GetRPCClient().NoteClient.Modify(context.Background(), &rpc_impl.Note_ModifyReq{
		OperatorId: params.OperatorID,
		NoteId:     params.NoteID,
		Password:   params.Password,
		Topic:      params.Topic,
		Content:    params.Content,
		IsPublic:   params.IsPublic,
	})
	if err != nil {
		mlog.Logger().Error("modify note failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func DeleteNote(r *http.Request) *mresponse.ResponseData {
	params := &structure.DeleteNoteReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	_, err := rpc.GetRPCClient().NoteClient.Delete(context.Background(), &rpc_impl.Note_DeleteReq{
		OperatorId: params.OperatorID,
		Password:   params.Password,
		NoteId:     params.NoteID,
	})
	if err != nil {
		mlog.Logger().Error("delete note failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func notesRPCToHTTP(data ...*rpc_impl.Note_Data) []*structure.NoteRes {
	res := make([]*structure.NoteRes, 0, len(data))
	for i := range data {
		res = append(res, &structure.NoteRes{
			NoteID:      data[i].NoteId,
			WriteBy:     data[i].WriteBy,
			Topic:       data[i].Topic,
			Content:     data[i].Content,
			IsPublic:    data[i].IsPublic,
			UpdateTime:  time.Duration(data[i].UpdateTime),
			CreatedTime: time.Duration(data[i].CreatedTime),
		})
	}

	return res
}
