package handlers

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/gateway/http/structure_defination"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/http"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/registration_center_embedded"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func getNoteClientAndConnTarget() (rpc_impl.INoteClient, string, error) {
	conn, err := rce.GetClientConn(mconst.UID_Service_Note)
	if err != nil {
		mlog.Logger().Error("get client conn failed", zap.Error(err))
		return nil, "", err
	}

	return rpc_impl.NewINoteClient(conn), conn.Target(), nil
}

func ListNoteByWriter(r *http.Request) *mhttp.ResponseData {
	params := &structure.ListNoteByWriterReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, target, err := getNoteClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get note client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.ListByWriter(context.Background(), &rpc_impl.Note_ListByWriterReq{
		OperatorId: params.OperatorID,
		Page: &rpc_impl.Pagination{
			PageSize: uint32(params.PageSize),
			PageNum:  uint32(params.PageNum),
		},
	})
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_Note, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(structure.MakeListNoteByWriterRes(res.Total, notesRPCToHTTP(res.Notes...)))
}

func ListPublicNote(r *http.Request) *mhttp.ResponseData {
	params := &structure.ListPublicNoteReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, target, err := getNoteClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get note client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.ListPublic(context.Background(), &rpc_impl.Note_ListPublicReq{
		OperatorId: params.OperatorID,
		Page: &rpc_impl.Pagination{
			PageSize: uint32(params.PageSize),
			PageNum:  uint32(params.PageNum),
		},
	})
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_Note, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(structure.MakeListPublicNoteRes(res.Total, notesRPCToHTTP(res.Notes...)))
}

func CreateNote(r *http.Request) *mhttp.ResponseData {
	params := &structure.CreateNoteReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, target, err := getNoteClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get note client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.Create(context.Background(), &rpc_impl.Note_CreateReq{
		OperatorId: params.OperatorID,
		Topic:      params.Topic,
		Content:    params.Content,
		IsPublic:   params.IsPublic,
	})
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_Note, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func ModifyNote(r *http.Request) *mhttp.ResponseData {
	params := &structure.ModifyNoteReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, target, err := getNoteClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get note client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.Modify(context.Background(), &rpc_impl.Note_ModifyReq{
		OperatorId: params.OperatorID,
		NoteId:     params.NoteID,
		Password:   params.Password,
		Topic:      params.Topic,
		Content:    params.Content,
		IsPublic:   params.IsPublic,
	})
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_Note, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func DeleteNote(r *http.Request) *mhttp.ResponseData {
	params := &structure.DeleteNoteReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, target, err := getNoteClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get note client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.Delete(context.Background(), &rpc_impl.Note_DeleteReq{
		OperatorId: params.OperatorID,
		Password:   params.Password,
		NoteId:     params.NoteID,
	})
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_Note, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.ResponseWithError(res.Err.String())
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
