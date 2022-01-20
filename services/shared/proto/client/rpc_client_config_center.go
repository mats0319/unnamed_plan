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

func (c configCenterClient) GetConfig(ctx context.Context, in *rpc_impl.Config_GetConfigReq, opts ...grpc.CallOption) (*rpc_impl.Config_GetConfigRes, error) {
    return c.client.GetConfig(ctx, in, opts...)
}
