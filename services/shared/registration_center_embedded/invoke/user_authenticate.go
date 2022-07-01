package rce_invoke

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/registration_center_embedded"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
)

func AuthUserInfo(ctx context.Context, userID string, password string) *rpc_impl.Error {
	conn, err := rce.GetClientConn(mconst.UID_Service_User)
	if err != nil {
		mlog.Logger().Error("get client conn failed", zap.Error(err))
		return utils.NewExecError(err.Error()).ToRPC()
	}

	authRes, err := rpc_impl.NewIUserClient(conn).Authenticate(ctx, &rpc_impl.User_AuthenticateReq{
		UserId:   userID,
		Password: password,
	})
	if err != nil { // grpc connection error
		rce.ReportInvalidTarget(mconst.UID_Service_User, conn.Target())

		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))

		return utils.NewGrpcConnectionError(err.Error()).ToRPC()
	}
	if authRes != nil && authRes.Err != nil { // service exec error
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", authRes.Err.String()))

		return authRes.Err
	}

	return nil
}
