package client

import (
    "context"
    "github.com/mats9693/unnamed_plan/services/shared/proto/impl"
    "google.golang.org/grpc"
)

type thinkingNoteClient struct {
    client rpc_impl.IThinkingNoteClient
}

var _ rpc_impl.IThinkingNoteClient = (*thinkingNoteClient)(nil)

var thinkingNoteClientIns = &thinkingNoteClient{}

func ConnectThinkingNoteServer(target string) (*thinkingNoteClient, error) {
    conn, err := grpc.Dial(target, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }

    thinkingNoteClientIns.client = rpc_impl.NewIThinkingNoteClient(conn)

    return thinkingNoteClientIns, nil
}

func (t *thinkingNoteClient) ListByWriter(ctx context.Context, in *rpc_impl.ThinkingNote_ListByWriterReq, opts ...grpc.CallOption) (*rpc_impl.ThinkingNote_ListByWriterRes, error) {
    return t.client.ListByWriter(ctx, in, opts...)
}

func (t *thinkingNoteClient) ListPublic(ctx context.Context, in *rpc_impl.ThinkingNote_ListPublicReq, opts ...grpc.CallOption) (*rpc_impl.ThinkingNote_ListPublicRes, error) {
    return t.client.ListPublic(ctx, in, opts...)
}

func (t *thinkingNoteClient) Create(ctx context.Context, in *rpc_impl.ThinkingNote_CreateReq, opts ...grpc.CallOption) (*rpc_impl.ThinkingNote_CreateRes, error) {
    return t.client.Create(ctx, in, opts...)
}

func (t *thinkingNoteClient) Modify(ctx context.Context, in *rpc_impl.ThinkingNote_ModifyReq, opts ...grpc.CallOption) (*rpc_impl.ThinkingNote_ModifyRes, error) {
    return t.client.Modify(ctx, in, opts...)
}

func (t *thinkingNoteClient) Delete(ctx context.Context, in *rpc_impl.ThinkingNote_DeleteReq, opts ...grpc.CallOption) (*rpc_impl.ThinkingNote_DeleteRes, error) {
    return t.client.Delete(ctx, in, opts...)
}
