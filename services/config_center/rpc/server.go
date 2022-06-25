package rpc

import (
	"context"
	"encoding/json"
	i "github.com/mats9693/unnamed_plan/services/config_center/init"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
)

type configCenterServerImpl struct {
	registrationCenterCoreTarget string

	rpc_impl.UnimplementedIConfigCenterServer
	rpc_impl.UnimplementedIConfigCenterRCServer
}

// auto-generate method declaration of service by IDE
//var _ rpc_impl.IConfigCenterServer = (*configCenterServerImpl)(nil)
//var _ rpc_impl.IConfigCenterRCServer = (*configCenterServerImpl)(nil)

var configCenterServerImplIns = &configCenterServerImpl{}

func GetConfigCenterServer() *configCenterServerImpl {
	return configCenterServerImplIns
}

func (c *configCenterServerImpl) GetServiceConfig(_ context.Context, req *rpc_impl.ConfigCenter_GetServiceConfigReq) (*rpc_impl.ConfigCenter_GetServiceConfigRes, error) {
	res := &rpc_impl.ConfigCenter_GetServiceConfigRes{}

	if len(req.ServiceId) < 1 || len(req.Level) < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("service id", req.ServiceId),
			zap.String("level", req.Level))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	config := i.GetServiceConfig(req.ServiceId, req.Level)
	if config == nil {
		mlog.Logger().Error("get service config failed",
			zap.String("service id", req.ServiceId),
			zap.String("level", req.Level))
		res.Err = utils.NewExecError("get service config failed").ToRPC()
		return res, nil
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		mlog.Logger().Error("json marshal failed", zap.Error(err))
		res.Err = utils.NewExecError("json marshal failed").ToRPC()
		return res, nil
	}

	res.Config = string(configBytes)
	res.RcCoreTarget = c.registrationCenterCoreTarget

	return res, nil
}

func (c *configCenterServerImpl) SetRCCoreTarget(_ context.Context, req *rpc_impl.ConfigCenterRC_SetRCCoreTargetReq) (*rpc_impl.ConfigCenterRC_SetRCCoreTargetRes, error) {
	res := &rpc_impl.ConfigCenterRC_SetRCCoreTargetRes{}

	if len(req.Target) < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams, zap.String("target", req.Target))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	c.registrationCenterCoreTarget = req.Target

	return res, nil
}
