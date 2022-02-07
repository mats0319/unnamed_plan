package client

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"google.golang.org/grpc"
)

type taskClient struct {
	client rpc_impl.ITaskClient
}

var _ rpc_impl.ITaskClient = (*taskClient)(nil)

var taskClientIns = &taskClient{}

func ConnectTaskServer(target string) (*taskClient, error) {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	taskClientIns.client = rpc_impl.NewITaskClient(conn)

	return taskClientIns, nil
}

func (t *taskClient) List(ctx context.Context, in *rpc_impl.Task_ListReq, opts ...grpc.CallOption) (*rpc_impl.Task_ListRes, error) {
	return t.client.List(ctx, in, opts...)
}

func (t *taskClient) Create(ctx context.Context, in *rpc_impl.Task_CreateReq, opts ...grpc.CallOption) (*rpc_impl.Task_CreateRes, error) {
	return t.client.Create(ctx, in, opts...)
}

func (t *taskClient) Modify(ctx context.Context, in *rpc_impl.Task_ModifyReq, opts ...grpc.CallOption) (*rpc_impl.Task_ModifyRes, error) {
	return t.client.Modify(ctx, in, opts...)
}
