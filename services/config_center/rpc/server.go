package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	i "github.com/mats9693/unnamed_plan/services/config_center/init"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
)

type configCenterServerImpl struct {
	rpc_impl.UnimplementedIConfigCenterServer
}

var _ rpc_impl.IConfigCenterServer = (*configCenterServerImpl)(nil)

var configCenterServerImplIns = &configCenterServerImpl{}

func GetConfigCenterServer() *configCenterServerImpl {
	return configCenterServerImplIns
}

func (c *configCenterServerImpl) GetServiceConfig(_ context.Context, req *rpc_impl.ConfigCenter_GetServiceConfigReq) (*rpc_impl.ConfigCenter_GetServiceConfigRes, error) {
	if len(req.ServiceId) < 1 || len(req.Level) < 1 {
		return nil, utils.NewError(mconst.Error_InvalidParams)
	}

	config := i.GetServiceConfig(req.ServiceId, req.Level)
	if config == nil {
		return nil, utils.NewError(fmt.Sprintf("unsupported config, service id: %s, level: %s\n",
			req.ServiceId, req.Level))
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return &rpc_impl.ConfigCenter_GetServiceConfigRes{
		Config: string(configBytes),
	}, nil
}
