package client

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"google.golang.org/grpc"
)

type noteClient struct {
	client rpc_impl.INoteClient
}

var _ rpc_impl.INoteClient = (*noteClient)(nil)

var thinkingNoteClientIns = &noteClient{}

func ConnectNoteServer(target string) (*noteClient, error) {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	thinkingNoteClientIns.client = rpc_impl.NewINoteClient(conn)

	return thinkingNoteClientIns, nil
}

func (t *noteClient) ListByWriter(ctx context.Context, in *rpc_impl.Note_ListByWriterReq, opts ...grpc.CallOption) (*rpc_impl.Note_ListByWriterRes, error) {
	return t.client.ListByWriter(ctx, in, opts...)
}

func (t *noteClient) ListPublic(ctx context.Context, in *rpc_impl.Note_ListPublicReq, opts ...grpc.CallOption) (*rpc_impl.Note_ListPublicRes, error) {
	return t.client.ListPublic(ctx, in, opts...)
}

func (t *noteClient) Create(ctx context.Context, in *rpc_impl.Note_CreateReq, opts ...grpc.CallOption) (*rpc_impl.Note_CreateRes, error) {
	return t.client.Create(ctx, in, opts...)
}

func (t *noteClient) Modify(ctx context.Context, in *rpc_impl.Note_ModifyReq, opts ...grpc.CallOption) (*rpc_impl.Note_ModifyRes, error) {
	return t.client.Modify(ctx, in, opts...)
}

func (t *noteClient) Delete(ctx context.Context, in *rpc_impl.Note_DeleteReq, opts ...grpc.CallOption) (*rpc_impl.Note_DeleteRes, error) {
	return t.client.Delete(ctx, in, opts...)
}
