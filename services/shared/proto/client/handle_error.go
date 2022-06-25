package client

import (
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/registration_center_embedded"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
)

func handleError(err error, serviceID string, target string, rpcErr *rpc_impl.Error) *rpc_impl.Error {
	if err != nil { // grpc connection error
		rc_embedded.ReportInvalidTarget(serviceID, target)

		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))

		return utils.NewGrpcConnectionError(err.Error()).ToRPC()
	}

	if rpcErr != nil { // service exec error
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", rpcErr.String()))

		return utils.NewExecError(rpcErr.String()).ToRPC()
	}

	return nil
}
