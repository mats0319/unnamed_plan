package client

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"google.golang.org/grpc"
)

type userClient struct {
	client rpc_impl.IUserClient
}

var _ rpc_impl.IUserClient = (*userClient)(nil)

var userClientIns = &userClient{}

func ConnectUserServer(target string) (*userClient, error) {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	userClientIns.client = rpc_impl.NewIUserClient(conn)

	return userClientIns, nil
}

func (u *userClient) Login(ctx context.Context, in *rpc_impl.User_LoginReq, opts ...grpc.CallOption) (*rpc_impl.User_LoginRes, error) {
	return u.client.Login(ctx, in, opts...)
}

func (u *userClient) List(ctx context.Context, in *rpc_impl.User_ListReq, opts ...grpc.CallOption) (*rpc_impl.User_ListRes, error) {
	return u.client.List(ctx, in, opts...)
}

func (u *userClient) Create(ctx context.Context, in *rpc_impl.User_CreateReq, opts ...grpc.CallOption) (*rpc_impl.User_CreateRes, error) {
	return u.client.Create(ctx, in, opts...)
}

func (u *userClient) Lock(ctx context.Context, in *rpc_impl.User_LockReq, opts ...grpc.CallOption) (*rpc_impl.User_LockRes, error) {
	return u.client.Lock(ctx, in, opts...)
}

func (u *userClient) Unlock(ctx context.Context, in *rpc_impl.User_UnlockReq, opts ...grpc.CallOption) (*rpc_impl.User_UnlockRes, error) {
	return u.client.Unlock(ctx, in, opts...)
}

func (u *userClient) ModifyInfo(ctx context.Context, in *rpc_impl.User_ModifyInfoReq, opts ...grpc.CallOption) (*rpc_impl.User_ModifyInfoRes, error) {
	return u.client.ModifyInfo(ctx, in, opts...)
}

func (u *userClient) ModifyPermission(ctx context.Context, in *rpc_impl.User_ModifyPermissionReq, opts ...grpc.CallOption) (*rpc_impl.User_ModifyPermissionRes, error) {
	return u.client.ModifyPermission(ctx, in, opts...)
}

func (u *userClient) Authenticate(ctx context.Context, in *rpc_impl.User_AuthenticateReq, opts ...grpc.CallOption) (*rpc_impl.User_AuthenticateRes, error) {
	return u.client.Authenticate(ctx, in, opts...)
}
