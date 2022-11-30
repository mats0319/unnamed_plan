package rpc

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/go"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func checkHealth(rceTarget string) (isHealthy bool) {
	targetSlice, err := utils.FormatTarget(rceTarget)
	if err == nil && len(targetSlice) > 0 {
		rceTarget = targetSlice[0]
	}

	conn, err := grpc.Dial(rceTarget, grpc.WithInsecure())
	if err != nil {
		mlog.Logger().Error("grpc connection error", zap.Error(err))
		return
	}

	client := rpc_impl.NewIRegistrationCenterEmbeddedClient(conn)

	res, err := client.CheckHealth(context.Background(), &rpc_impl.RegistrationCenterEmbedded_CheckHealthReq{
		Data: "", // ^_^
	})
	if err != nil {
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		return
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return
	}

	isHealthy = true

	return
}
