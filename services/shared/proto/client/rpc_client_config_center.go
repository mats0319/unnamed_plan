package client

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"google.golang.org/grpc"
)

type configCenterClient struct {
	client rpc_impl.IConfigCenterClient
}

var _ rpc_impl.IConfigCenterClient = (*configCenterClient)(nil)

var configCenterClientIns = &configCenterClient{}

func ConnectConfigCenterServer(target string) (*configCenterClient, error) {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	configCenterClientIns.client = rpc_impl.NewIConfigCenterClient(conn)

	return configCenterClientIns, nil
}

func (c *configCenterClient) GetServiceConfig(ctx context.Context, in *rpc_impl.ConfigCenter_GetServiceConfigReq, opts ...grpc.CallOption) (*rpc_impl.ConfigCenter_GetServiceConfigRes, error) {
	return c.client.GetServiceConfig(ctx, in, opts...)
}
