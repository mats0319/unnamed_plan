package rpc

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/thinking_note/db"
	"github.com/mats9693/unnamed_plan/shared/db/model"
	"github.com/mats9693/unnamed_plan/shared/errors"
	"github.com/mats9693/unnamed_plan/shared/proto/client"
	"github.com/mats9693/unnamed_plan/shared/proto/impl"
	"github.com/pkg/errors"
)

type thinkingNoteServerImpl struct {
	rpc_impl.UnimplementedIThinkingNoteServer

	UserClient rpc_impl.IUserClient
}

var thinkingNoteServerImplIns = &thinkingNoteServerImpl{}

var _ rpc_impl.IThinkingNoteServer = (*thinkingNoteServerImpl)(nil)

func GetThinkingNoteServer(userServerTarget string) (*thinkingNoteServerImpl, error) {
    userClient, err := client.ConnectUserServer(userServerTarget)
    if err != nil {
        return nil, err
    }

    thinkingNoteServerImplIns.UserClient = userClient

	return thinkingNoteServerImplIns,nil
}

func (t *thinkingNoteServerImpl) ListByWriter(_ context.Context, req *rpc_impl.ThinkingNote_ListByWriterReq) (*rpc_impl.ThinkingNote_ListByWriterRes, error) {
	if len(req.OperatorId) < 1 || req.GetPage() == nil || req.GetPage().PageSize < 1 || req.GetPage().PageNum < 1 {
		return nil, errors.New(merrors.Error_InvalidParams)
	}

	pageSize := int(req.GetPage().PageSize)
	pageNum := int(req.GetPage().PageNum)

	notes, count, err := db.GetThinkingNote().QueryPageByWriter(pageSize, pageNum, req.OperatorId)
	if err != nil {
		return nil, err
	}

	return &rpc_impl.ThinkingNote_ListByWriterRes{
		Total: uint32(count),
		Notes: notesDBToRPC(notes...),
	}, nil
}

func (t *thinkingNoteServerImpl) ListPublic(_ context.Context, req *rpc_impl.ThinkingNote_ListPublicReq) (*rpc_impl.ThinkingNote_ListPublicRes, error) {
	if len(req.OperatorId) < 1 || req.GetPage() == nil || req.GetPage().PageSize < 1 || req.GetPage().PageNum < 1 {
		return nil, errors.New(merrors.Error_InvalidParams)
	}

	pageSize := int(req.GetPage().PageSize)
	pageNum := int(req.GetPage().PageNum)

	notes, count, err := db.GetThinkingNote().QueryPageInPublic(pageSize, pageNum, req.OperatorId)
	if err != nil {
		return nil, err
	}

	return &rpc_impl.ThinkingNote_ListPublicRes{
		Total: uint32(count),
		Notes: notesDBToRPC(notes...),
	}, nil
}

func (t *thinkingNoteServerImpl) Create(_ context.Context, req *rpc_impl.ThinkingNote_CreateReq) (*rpc_impl.ThinkingNote_CreateRes, error) {
	if len(req.OperatorId) < 1 || len(req.Content) < 1 {
		return nil, errors.New(merrors.Error_InvalidParams)
	}

	err := db.GetThinkingNote().Insert(&model.ThinkingNote{
		WriteBy:  req.OperatorId,
		Topic:    req.Topic,
		Content:  req.Content,
		IsPublic: req.IsPublic,
	})
	if err != nil {
		return nil, err
	}

	return &rpc_impl.ThinkingNote_CreateRes{}, nil
}

func (t *thinkingNoteServerImpl) Modify(ctx context.Context, req *rpc_impl.ThinkingNote_ModifyReq) (*rpc_impl.ThinkingNote_ModifyRes, error) {
	if len(req.OperatorId) < 1 || len(req.NoteId) < 1 {
		return nil, errors.New(merrors.Error_InvalidParams)
	}

	_, err := t.UserClient.Authenticate(ctx, &rpc_impl.User_AuthenticateReq{
        UserId:   req.OperatorId,
        Password: req.Password,
    })
	if err != nil {
		return nil, err
	}

	noteRecord, err := db.GetThinkingNote().QueryFirst(model.ThinkingNote_NoteID+" = ?", req.NoteId)
	if err != nil {
		return nil, err
	}
	if noteRecord.WriteBy != req.OperatorId {
		return nil, errors.New(merrors.Error_ModifyOthersThinkingNote)
	}

	if len(req.Topic)+len(req.Content) < 1 && noteRecord.IsPublic == req.IsPublic {
		return nil, errors.New(merrors.Error_NoValidModification)
	}

	updateColumns := make([]string, 0, 3)
	if len(req.Topic) > 0 {
		noteRecord.Topic = req.Topic
		updateColumns = append(updateColumns, model.ThinkingNote_Topic)
	}
	if len(req.Content) > 0 {
		noteRecord.Content = req.Content
		updateColumns = append(updateColumns, model.ThinkingNote_Content)
	}
	if noteRecord.IsPublic != req.IsPublic {
		noteRecord.IsPublic = req.IsPublic
		updateColumns = append(updateColumns, model.ThinkingNote_IsPublic)
	}

	err = db.GetThinkingNote().UpdateColumnsByNoteID(noteRecord, updateColumns...)
	if err != nil {
		return nil, err
	}

	return &rpc_impl.ThinkingNote_ModifyRes{}, nil
}

func (t *thinkingNoteServerImpl) Delete(ctx context.Context, req *rpc_impl.ThinkingNote_DeleteReq) (*rpc_impl.ThinkingNote_DeleteRes, error) {
    if len(req.OperatorId) < 1 || len(req.Password) < 1 || len(req.NoteId) < 1 {
        return nil, errors.New(merrors.Error_InvalidParams)
    }

	_, err := t.UserClient.Authenticate(ctx, &rpc_impl.User_AuthenticateReq{
		UserId:   req.OperatorId,
		Password: req.Password,
	})
    if err != nil {
        return nil, err
    }

    err = db.GetThinkingNote().UpdateColumnsByNoteID(&model.ThinkingNote{
        NoteID:    req.NoteId,
        IsDeleted: true,
    }, model.ThinkingNote_IsDeleted)
    if err != nil {
        return nil, err
    }

	return &rpc_impl.ThinkingNote_DeleteRes{}, nil
}

func notesDBToRPC(data ...*model.ThinkingNote) []*rpc_impl.ThinkingNote_Data {
	res := make([]*rpc_impl.ThinkingNote_Data, 0, len(data))
	for i := range data {
		res = append(res, &rpc_impl.ThinkingNote_Data{
			NoteId:      data[i].NoteID,
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
