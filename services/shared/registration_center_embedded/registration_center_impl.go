package rc_embedded

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
)

// rcEmbeddedImpl communicate with RC core
type rcEmbeddedImpl struct {
	rcCoreClient rpc_impl.IRegistrationCenterCoreClient

	// todo: rc embedded service impl
}

// register make sure 'r' is initialized before invoke
func (r *rcEmbeddedImpl) register(serviceID string, target string) error {
	res, err := r.rcCoreClient.Register(context.Background(), &rpc_impl.RegistrationCenterCore_RegisterReq{
		ServiceId: serviceID,
		Target:    target,
	})
	if err != nil || (res == nil || res.Err != nil) {
		mlog.Logger().Error("register service failed",
			zap.String("service id", serviceID),
			zap.String("target", target),
			zap.NamedError("error", err),
			zap.Any("res", res))
		return utils.NewError(err.Error() + res.Err.String())
	}

	return nil
}

// ListServiceTarget make sure 'r' is initialized before invoke
func (r *rcEmbeddedImpl) ListServiceTarget(serviceID string) ([]string, error) {
	res, err := r.rcCoreClient.ListServiceTarget(context.Background(), &rpc_impl.RegistrationCenterCore_ListServiceTargetReq{
		ServiceId: serviceID,
	})
	if err != nil || (res == nil || res.Err != nil) {
		mlog.Logger().Error("list service target failed",
			zap.String("service id", serviceID),
            zap.NamedError("error", err),
			zap.Any("res", res))
		return nil, utils.NewError(err.Error() + res.Err.String())
	}

	return res.Targets, nil
}
