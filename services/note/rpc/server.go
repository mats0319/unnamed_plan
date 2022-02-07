package rpc

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/note/db"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"github.com/mats9693/unnamed_plan/services/shared/proto/client"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
)

type noteServerImpl struct {
	rpc_impl.UnimplementedINoteServer

	UserClient rpc_impl.IUserClient
}

var thinkingNoteServerImplIns = &noteServerImpl{}

var _ rpc_impl.INoteServer = (*noteServerImpl)(nil)

func GetNoteServer(userServerTarget string) (*noteServerImpl, error) {
	userClient, err := client.ConnectUserServer(userServerTarget)
	if err != nil {
		return nil, err
	}

	thinkingNoteServerImplIns.UserClient = userClient

	return thinkingNoteServerImplIns, nil
}

func (t *noteServerImpl) ListByWriter(_ context.Context, req *rpc_impl.Note_ListByWriterReq) (*rpc_impl.Note_ListByWriterRes, error) {
	if len(req.OperatorId) < 1 || req.GetPage() == nil || req.GetPage().PageSize < 1 || req.GetPage().PageNum < 1 {
		return nil, utils.NewError(mconst.Error_InvalidParams)
	}

	pageSize := int(req.GetPage().PageSize)
	pageNum := int(req.GetPage().PageNum)

	notes, count, err := db.GetNoteDao().QueryPageByWriter(pageSize, pageNum, req.OperatorId)
	if err != nil {
		return nil, err
	}

	return &rpc_impl.Note_ListByWriterRes{
		Total: uint32(count),
		Notes: notesDBToRPC(notes...),
	}, nil
}

func (t *noteServerImpl) ListPublic(_ context.Context, req *rpc_impl.Note_ListPublicReq) (*rpc_impl.Note_ListPublicRes, error) {
	if len(req.OperatorId) < 1 || req.GetPage() == nil || req.GetPage().PageSize < 1 || req.GetPage().PageNum < 1 {
		return nil, utils.NewError(mconst.Error_InvalidParams)
	}

	pageSize := int(req.GetPage().PageSize)
	pageNum := int(req.GetPage().PageNum)

	notes, count, err := db.GetNoteDao().QueryPageInPublic(pageSize, pageNum, req.OperatorId)
	if err != nil {
		return nil, err
	}

	return &rpc_impl.Note_ListPublicRes{
		Total: uint32(count),
		Notes: notesDBToRPC(notes...),
	}, nil
}

func (t *noteServerImpl) Create(_ context.Context, req *rpc_impl.Note_CreateReq) (*rpc_impl.Note_CreateRes, error) {
	if len(req.OperatorId) < 1 || len(req.Content) < 1 {
		return nil, utils.NewError(mconst.Error_InvalidParams)
	}

	err := db.GetNoteDao().Insert(&model.Note{
		WriteBy:  req.OperatorId,
		Topic:    req.Topic,
		Content:  req.Content,
		IsPublic: req.IsPublic,
	})
	if err != nil {
		return nil, err
	}

	return &rpc_impl.Note_CreateRes{}, nil
}

func (t *noteServerImpl) Modify(ctx context.Context, req *rpc_impl.Note_ModifyReq) (*rpc_impl.Note_ModifyRes, error) {
	if len(req.OperatorId) < 1 || len(req.NoteId) < 1 || len(req.Password) < 1 || len(req.Content) < 1 {
		return nil, utils.NewError(mconst.Error_InvalidParams)
	}

	_, err := t.UserClient.Authenticate(ctx, &rpc_impl.User_AuthenticateReq{
		UserId:   req.OperatorId,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	noteRecord, err := db.GetNoteDao().QueryOne(req.NoteId)
	if err != nil {
		return nil, err
	}
	if noteRecord.WriteBy != req.OperatorId {
		return nil, utils.NewError(mconst.Error_ModifyOthersThinkingNote)
	}

	if req.Topic == noteRecord.Topic && req.Content == noteRecord.Content && req.IsPublic == noteRecord.IsPublic {
		return nil, utils.NewError(mconst.Error_NoValidModification)
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
		return nil, err
	}

	return &rpc_impl.Note_ModifyRes{}, nil
}

func (t *noteServerImpl) Delete(ctx context.Context, req *rpc_impl.Note_DeleteReq) (*rpc_impl.Note_DeleteRes, error) {
	if len(req.OperatorId) < 1 || len(req.Password) < 1 || len(req.NoteId) < 1 {
		return nil, utils.NewError(mconst.Error_InvalidParams)
	}

	_, err := t.UserClient.Authenticate(ctx, &rpc_impl.User_AuthenticateReq{
		UserId:   req.OperatorId,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	note := &model.Note{}
	note.ID = req.NoteId
	note.IsDeleted = true

	err = db.GetNoteDao().UpdateColumnsByNoteID(note, model.Note_IsDeleted)
	if err != nil {
		return nil, err
	}

	return &rpc_impl.Note_DeleteRes{}, nil
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
