package rpc

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/3_note/db"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/client"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
)

type noteServerImpl struct {
	rpc_impl.UnimplementedINoteServer
}

var noteServerImplIns = &noteServerImpl{}

func GetNoteServer() rpc_impl.INoteServer {
	return noteServerImplIns
}

func (t *noteServerImpl) ListByWriter(_ context.Context, req *rpc_impl.Note_ListByWriterReq) (*rpc_impl.Note_ListByWriterRes, error) {
	res := &rpc_impl.Note_ListByWriterRes{}

	if len(req.OperatorId) < 1 || req.GetPage() == nil || req.Page.PageSize < 1 || req.Page.PageNum < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("operator", req.OperatorId),
			zap.String("page info", req.Page.String()))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	pageSize := int(req.Page.PageSize)
	pageNum := int(req.Page.PageNum)

	notes, count, err := db.GetNoteDao().QueryPageByWriter(pageSize, pageNum, req.OperatorId)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	res.Total = uint32(count)
	res.Notes = notesDBToRPC(notes...)

	return res, nil
}

func (t *noteServerImpl) ListPublic(_ context.Context, req *rpc_impl.Note_ListPublicReq) (*rpc_impl.Note_ListPublicRes, error) {
	res := &rpc_impl.Note_ListPublicRes{}

	if len(req.OperatorId) < 1 || req.GetPage() == nil || req.Page.PageSize < 1 || req.Page.PageNum < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("operator", req.OperatorId),
			zap.String("page info", req.Page.String()))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	pageSize := int(req.Page.PageSize)
	pageNum := int(req.Page.PageNum)

	notes, count, err := db.GetNoteDao().QueryPageInPublic(pageSize, pageNum, req.OperatorId)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	res.Total = uint32(count)
	res.Notes = notesDBToRPC(notes...)

	return res, nil
}

func (t *noteServerImpl) Create(_ context.Context, req *rpc_impl.Note_CreateReq) (*rpc_impl.Note_CreateRes, error) {
	res := &rpc_impl.Note_CreateRes{}

	if len(req.OperatorId) < 1 || len(req.Content) < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("operator", req.OperatorId),
			zap.String("content", req.Content))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	err := db.GetNoteDao().Insert(&model.Note{
		WriteBy:  req.OperatorId,
		Topic:    req.Topic,
		Content:  req.Content,
		IsPublic: req.IsPublic,
	})
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	return res, nil
}

func (t *noteServerImpl) Modify(ctx context.Context, req *rpc_impl.Note_ModifyReq) (*rpc_impl.Note_ModifyRes, error) {
	res := &rpc_impl.Note_ModifyRes{}

	if len(req.OperatorId) < 1 || len(req.NoteId) < 1 || len(req.Password) < 1 || len(req.Content) < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("operator", req.OperatorId),
			zap.String("note", req.NoteId),
			zap.String("password", req.Password),
			zap.String("content", req.Content))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	rpcErr := client.AuthUserInfo(ctx, req.OperatorId, req.Password)
	if rpcErr != nil {
		mlog.Logger().Error("auth user info failed", zap.String("error", rpcErr.String()))
		res.Err = rpcErr
		return res, nil
	}

	noteRecord, err := db.GetNoteDao().QueryOne(req.NoteId)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	if noteRecord.WriteBy != req.OperatorId {
		mlog.Logger().Error(mconst.Error_ModifyOthersNote)
		res.Err = utils.NewExecError(mconst.Error_ModifyOthersNote).ToRPC()
		return res, nil
	}

	if req.Topic == noteRecord.Topic && req.Content == noteRecord.Content && req.IsPublic == noteRecord.IsPublic {
		mlog.Logger().Error(mconst.Error_NoValidModification)
		res.Err = utils.NewExecError(mconst.Error_NoValidModification).ToRPC()
		return res, nil
	}

	updateColumns := make([]string, 0, 3)

	noteRecord.Content = req.Content
	updateColumns = append(updateColumns, model.Note_Content)

	if len(req.Topic) > 0 {
		noteRecord.Topic = req.Topic
		updateColumns = append(updateColumns, model.Note_Topic)
	}
	if noteRecord.IsPublic != req.IsPublic {
		noteRecord.IsPublic = req.IsPublic
		updateColumns = append(updateColumns, model.Note_IsPublic)
	}

	err = db.GetNoteDao().UpdateColumnsByNoteID(noteRecord, updateColumns...)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	return res, nil
}

func (t *noteServerImpl) Delete(ctx context.Context, req *rpc_impl.Note_DeleteReq) (*rpc_impl.Note_DeleteRes, error) {
	res := &rpc_impl.Note_DeleteRes{}

	if len(req.OperatorId) < 1 || len(req.Password) < 1 || len(req.NoteId) < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("operator", req.OperatorId),
			zap.String("password", req.Password),
			zap.String("note", req.NoteId))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	rpcErr := client.AuthUserInfo(ctx, req.OperatorId, req.Password)
	if rpcErr != nil {
		mlog.Logger().Error("auth user info failed", zap.String("error", rpcErr.String()))
		res.Err = rpcErr
		return res, nil
	}

	noteRecord, err := db.GetNoteDao().QueryOne(req.NoteId)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	noteRecord.IsDeleted = true

	err = db.GetNoteDao().UpdateColumnsByNoteID(noteRecord, model.Note_IsDeleted)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	return res, nil
}

func notesDBToRPC(data ...*model.Note) []*rpc_impl.Note_Data {
	res := make([]*rpc_impl.Note_Data, 0, len(data))
	for i := range data {
		res = append(res, &rpc_impl.Note_Data{
			NoteId:      data[i].ID,
			WriteBy:     data[i].WriteBy,
			Topic:       data[i].Topic,
			Content:     data[i].Content,
			IsPublic:    data[i].IsPublic,
			UpdateTime:  int64(data[i].UpdateTime),
			CreatedTime: int64(data[i].CreatedTime),
		})
	}

	return res
}
