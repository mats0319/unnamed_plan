package rce

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/go"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// rceImpl communicate with RC core
type rceImpl struct {
	init bool

	client rpc_impl.IRegistrationCenterCoreClient
	server rpc_impl.IRegistrationCenterEmbeddedServer
	rpc_impl.UnimplementedIRegistrationCenterEmbeddedServer
}

func (r *rceImpl) CheckHealth(_ context.Context, _ *rpc_impl.RegistrationCenterEmbedded_CheckHealthReq) (*rpc_impl.RegistrationCenterEmbedded_CheckHealthRes, error) {
	return &rpc_impl.RegistrationCenterEmbedded_CheckHealthRes{}, nil
}

// ListServiceTarget make sure 'r' is initialized before invoke
func (r *rceImpl) ListServiceTarget(serviceID string) ([]string, error) {
	if !r.init {
		return nil, errors.New("please init first")
	}

	res, err := r.client.ListServiceTarget(context.Background(), &rpc_impl.RegistrationCenterCore_ListServiceTargetReq{
		ServiceId: serviceID,
	})
	if err != nil {
		mlog.Logger().Error("list service target failed, with grpc connection error", zap.Error(err))
		return nil, err
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error("list service target failed", zap.String("error", res.Err.String()))
		return nil, utils.NewError(res.Err.String())
	}

	target, err := utils.FormatTarget(res.Targets...)
	if err != nil {
		mlog.Logger().Error("format target failed", zap.Error(err))
		return nil, err
	}

	return target, nil
}

func (r *rceImpl) initialize(target string) error {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		mlog.Logger().Error("grpc dial failed", zap.Error(err))
		return err
	}

	r.client = rpc_impl.NewIRegistrationCenterCoreClient(conn)

	r.init = true

	return nil
}

func (r *rceImpl) register(serviceID string, target string) error {
	if !r.init {
		mlog.Logger().Error("RCE module not init")
		return errors.New("uninitialized RCE module")
	}

	res, err := r.client.Register(context.Background(), &rpc_impl.RegistrationCenterCore_RegisterReq{
		ServiceId: serviceID,
		Target:    target,
	})
	if err != nil {
		mlog.Logger().Error("register service failed, with grpc connection error", zap.Error(err))
		return err
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error("register service failed", zap.String("error", res.Err.String()))
		return utils.NewError(res.Err.String())
	}

	return nil
}
