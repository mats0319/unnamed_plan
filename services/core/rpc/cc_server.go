package rpc

import (
	"context"
	"encoding/json"
	i "github.com/mats9693/unnamed_plan/services/core/init"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/go"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
)

type configCenterServerImpl struct {
	rpc_impl.UnimplementedIConfigCenterServer
}

var configCenterServerImplIns = &configCenterServerImpl{}

func GetConfigCenterServer() rpc_impl.IConfigCenterServer {
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

	return res, nil
}
