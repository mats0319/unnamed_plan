package rpc

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/go"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
	"sync"
	"time"
)

type registrationCenterServerImpl struct {
	targetMap sync.Map // service id - target slice([]string)

	rpc_impl.UnimplementedIRegistrationCenterCoreServer
}

var _ rpc_impl.IRegistrationCenterCoreServer = (*registrationCenterServerImpl)(nil)

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
		res.Err = utils.Error_InvalidParams
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
		res.Err = utils.Error_InvalidParams
		return res, nil
	}

	var targetSlice []string
	{
		targetSliceI, ok := r.targetMap.Load(req.ServiceId)
		targetSlice, _ = targetSliceI.([]string)
		if !ok || len(targetSlice) < 1 {
			mlog.Logger().Error("no valid service instance", zap.String("service id", req.ServiceId))
			targetSlice = make([]string, 0)
		}
	}

	res.Targets = targetSlice

	return res, nil
}

func (r *registrationCenterServerImpl) MaintainTarget(exitChan chan struct{}) {
	for {
		select {
		case <-exitChan:
			mlog.Logger().Error("stop maintain target, exit")
			return
		case <-time.After(time.Second * 60):
			mlog.Logger().Info("routine maintain targets")

			var (
				serviceIDs = make([]string, 0)
				targetList = make([][]string, 0)
			)

			// till this version, still can modify map in Range function, because we only modify value of current key
			// but consider for future features, do not modify map in Range function.
			// further more information, in Range function annotation
			//
			// get keys(service id) and values(target slice) of map
			r.targetMap.Range(func(keyI interface{}, valueI interface{}) bool {
				key, _ := keyI.(string)
				serviceIDs = append(serviceIDs, key)

				target, _ := valueI.([]string)
				targetList = append(targetList, target)

				return true // always return true, otherwise range will stop
			})

			if len(serviceIDs) != len(targetList) {
				mlog.Logger().Error("unexpected error, unmatched data length")
				break
			}

			// check health of each target, handle invalid target(s)
			//
			// in this version, delete invalid target(s) directly
			for i := range serviceIDs {
				newTarget := make([]string, 0, len(targetList[i]))

				for j := range targetList[i] {
					isHealthy := checkHealth(targetList[i][j])
					if isHealthy {
						newTarget = append(newTarget, targetList[i][j])
					} else {
						// delete invalid target
						mlog.Logger().Error("a service instance become unavailable",
							zap.String("service id", serviceIDs[i]),
							zap.String("target", targetList[i][j]))
					}
				}

				r.targetMap.Store(serviceIDs[i], newTarget)
			}
		}
	}
}
