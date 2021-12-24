package handlers

import (
    "context"
    "github.com/mats9693/unnamed_plan/services/gateway/http/structure_defination"
    "github.com/mats9693/unnamed_plan/services/gateway/rpc"
    "github.com/mats9693/unnamed_plan/services/shared/proto/impl"
    "github.com/mats9693/utils/toy_server/const"
    "github.com/mats9693/utils/toy_server/http"
    "github.com/mats9693/utils/toy_server/log"
    "go.uber.org/zap"
    "net/http"
    "time"
)

func ListThinkingNoteByWriter(r *http.Request) *mhttp.ResponseData {
    params := &structure.ListThinkingNoteByWriterReqParams{}
    if errMsg := params.Decode(r); len(errMsg) > 0 {
        mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
        return mhttp.ResponseWithError(errMsg)
    }

    res, err := rpc.GetRPCClient().ThinkingNoteClient.ListByWriter(context.Background(), &rpc_impl.ThinkingNote_ListByWriterReq{
        OperatorId: params.OperatorID,
        Page: &rpc_impl.Pagination{
            PageSize: uint32(params.PageSize),
            PageNum:  uint32(params.PageNum),
        },
    })
    if err != nil {
        mlog.Logger().Error("list thinking note by writer failed", zap.Error(err))
        return mhttp.ResponseWithError(err.Error())
    }

    return mhttp.Response(structure.MakeListThinkingNoteByWriterRes(res.Total, notesRPCToHTTP(res.Notes...)))
}

func ListPublicThinkingNote(r *http.Request) *mhttp.ResponseData {
    params := &structure.ListPublicThinkingNoteReqParams{}
    if errMsg := params.Decode(r); len(errMsg) > 0 {
        mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
        return mhttp.ResponseWithError(errMsg)
    }

    res, err := rpc.GetRPCClient().ThinkingNoteClient.ListPublic(context.Background(), &rpc_impl.ThinkingNote_ListPublicReq{
        OperatorId: params.OperatorID,
        Page: &rpc_impl.Pagination{
            PageSize: uint32(params.PageSize),
            PageNum:  uint32(params.PageNum),
        },
    })
    if err != nil {
        mlog.Logger().Error("list public thinking note failed", zap.Error(err))
        return mhttp.ResponseWithError(err.Error())
    }

    return mhttp.Response(structure.MakeListPublicThinkingNoteRes(res.Total, notesRPCToHTTP(res.Notes...)))
}

func CreateThinkingNote(r *http.Request) *mhttp.ResponseData {
    params := &structure.CreateThinkingNoteReqParams{}
    if errMsg := params.Decode(r); len(errMsg) > 0 {
        mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
        return mhttp.ResponseWithError(errMsg)
    }

    _, err := rpc.GetRPCClient().ThinkingNoteClient.Create(context.Background(), &rpc_impl.ThinkingNote_CreateReq{
        OperatorId: params.OperatorID,
        Topic:      params.Topic,
        Content:    params.Content,
        IsPublic:   params.IsPublic,
    })
    if err != nil {
        mlog.Logger().Error("create thinking note failed", zap.Error(err))
        return mhttp.ResponseWithError(err.Error())
    }

    return mhttp.Response(mconst.EmptyHTTPRes)
}

func ModifyThinkingNote(r *http.Request) *mhttp.ResponseData {
    params := &structure.ModifyThinkingNoteReqParams{}
    if errMsg := params.Decode(r); len(errMsg) > 0 {
        mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
        return mhttp.ResponseWithError(errMsg)
    }

    _, err := rpc.GetRPCClient().ThinkingNoteClient.Modify(context.Background(), &rpc_impl.ThinkingNote_ModifyReq{
        OperatorId: params.OperatorID,
        NoteId:     params.NoteID,
        Password:   params.Password,
        Topic:      params.Topic,
        Content:    params.Content,
        IsPublic:   params.IsPublic,
    })
    if err != nil {
        mlog.Logger().Error("modify thinking note failed", zap.Error(err))
        return mhttp.ResponseWithError(err.Error())
    }

    return mhttp.Response(mconst.EmptyHTTPRes)
}

func DeleteThinkingNote(r *http.Request) *mhttp.ResponseData {
    params := &structure.DeleteThinkingNoteReqParams{}
    if errMsg := params.Decode(r); len(errMsg) > 0 {
        mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
        return mhttp.ResponseWithError(errMsg)
    }

    _, err := rpc.GetRPCClient().ThinkingNoteClient.Delete(context.Background(), &rpc_impl.ThinkingNote_DeleteReq{
        OperatorId: params.OperatorID,
        Password:   params.Password,
        NoteId:     params.NoteID,
    })
    if err != nil {
        mlog.Logger().Error("delete thinking note failed", zap.Error(err))
        return mhttp.ResponseWithError(err.Error())
    }

    return mhttp.Response(mconst.EmptyHTTPRes)
}

func notesRPCToHTTP(data ...*rpc_impl.ThinkingNote_Data) []*structure.NoteRes {
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
