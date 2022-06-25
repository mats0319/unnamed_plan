package rpc

import (
    "context"
    "github.com/mats9693/unnamed_plan/services/shared/const"
    "github.com/mats9693/unnamed_plan/services/shared/log"
    "github.com/mats9693/unnamed_plan/services/shared/proto/impl"
    "github.com/mats9693/unnamed_plan/services/shared/utils"
    "go.uber.org/zap"
    "sync"
)

type registrationCenterServerImpl struct {
    targetMap sync.Map // service id - target slice([]string)

    rpc_impl.UnimplementedIRegistrationCenterCoreServer
}

// auto-generate method declaration of service by IDE
//var _ rpc_impl.IRegistrationCenterCoreServer = (*registrationCenterServerImpl)(nil)

var registrationCenterServerImplIns = &registrationCenterServerImpl{}

func GetRegistrationCenterServer() *registrationCenterServerImpl {
    return registrationCenterServerImplIns
}

func (r *registrationCenterServerImpl) Register(_ context.Context, req *rpc_impl.RegistrationCenterCore_RegisterReq) (*rpc_impl.RegistrationCenterCore_RegisterRes, error) {
    res := &rpc_impl.RegistrationCenterCore_RegisterRes{}

    if len(req.ServiceId) < 1 || len(req.Target) < 1 {
        mlog.Logger().Error(mconst.Error_InvalidParams,
            zap.String("service id", req.ServiceId),
            zap.String("target", req.Target))
        res.Err = utils.Error_InvalidParams.ToRPC()
        return res, nil
    }

    targetSliceI, ok := r.targetMap.Load(req.ServiceId)
    targetSlice, _ := targetSliceI.([]string)
    if !ok {
        targetSlice = []string{req.Target}
    } else {
        targetSlice = append(targetSlice, req.Target)
    }

    r.targetMap.Store(req.ServiceId, targetSlice)

    mlog.Logger().Info("a new service instance registered",
        zap.String("service id", req.ServiceId),
        zap.String("target", req.Target))

    return res, nil
}

func (r *registrationCenterServerImpl) ListServiceTarget(_ context.Context, req *rpc_impl.RegistrationCenterCore_ListServiceTargetReq) (*rpc_impl.RegistrationCenterCore_ListServiceTargetRes, error) {
    res := &rpc_impl.RegistrationCenterCore_ListServiceTargetRes{}

    if len(req.ServiceId) < 1 {
        mlog.Logger().Error(mconst.Error_InvalidParams, zap.String("service id", req.ServiceId))
        res.Err = utils.Error_InvalidParams.ToRPC()
        return res, nil
    }

    targetSliceI, ok := r.targetMap.Load(req.ServiceId)
    targetSlice, _ := targetSliceI.([]string)
    if !ok || len(targetSlice) < 1 {
        mlog.Logger().Error("no valid service instance", zap.String("service id", req.ServiceId))
        targetSlice = make([]string, 0)
    }

    res.Targets = targetSlice

    return res, nil
}
