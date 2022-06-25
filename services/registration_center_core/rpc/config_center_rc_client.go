package rpc

import (
    "context"
    "github.com/mats9693/unnamed_plan/services/shared/log"
    "github.com/mats9693/unnamed_plan/services/shared/proto/impl"
    "github.com/mats9693/unnamed_plan/services/shared/utils"
    "github.com/pkg/errors"
    "go.uber.org/zap"
    "google.golang.org/grpc"
)

type configCenterRCClientImpl struct {
    init bool

	client rpc_impl.IConfigCenterRCClient
}

var configCenterRCClientImplIns = &configCenterRCClientImpl{}

func GetConfigCenterRCClient() *configCenterRCClientImpl {
    return configCenterRCClientImplIns
}

func (r *configCenterRCClientImpl) Init(target string) error {
    conn, err := grpc.Dial(target, grpc.WithInsecure())
    if err != nil {
        mlog.Logger().Error("grpc dial failed", zap.Error(err))
        return err
    }

    r.client = rpc_impl.NewIConfigCenterRCClient(conn)

    r.init = true

    return nil
}

func (r *configCenterRCClientImpl) SetRCCoreTarget(target string) error {
    if !r.init {
        return errors.New("please init first")
    }

	res, err := r.client.SetRCCoreTarget(context.Background(), &rpc_impl.ConfigCenterRC_SetRCCoreTargetReq{
		Target: target,
	})
	if err != nil || (res == nil || res.Err != nil) {
        mlog.Logger().Error("set rc core target failed",
            zap.String("target", target),
            zap.NamedError("error", err),
            zap.Any("res", res))
        return utils.NewError(err.Error() + res.Err.String())
	}

    return nil
}
