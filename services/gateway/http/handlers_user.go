package http

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/gateway/http/structure"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/http"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/go"
	"github.com/mats9693/unnamed_plan/services/shared/registration_center_embedded"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
	"net/http"
)

// getUserClient may contain re-connect or other common handles in the future
func getUserClientAndConnTarget() (rpc_impl.IUserClient, string, error) {
	conn, err := rce.GetClientConn(mconst.UID_Service_User)
	if err != nil {
		mlog.Logger().Error("get client conn failed", zap.Error(err))
		return nil, "", err
	}

	return rpc_impl.NewIUserClient(conn), conn.Target(), nil
}

func Login(r *http.Request) *mhttp.ResponseData {
	res := &rpc_impl.User_LoginRes{}

	params := &structure.LoginReqParams{}
	params.Decode(r)

	client, target, err := getUserClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get user client failed", zap.Error(err))
		res.Err = utils.NewGetClientError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}

	res, err = client.Login(context.Background(), (*rpc_impl.User_LoginReq)(params))
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_User, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		res.Err = utils.NewGrpcConnectionError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.NewResponseDataWithError(res)
	}

	return mhttp.NewResponseData(res, res.UserId)
}

func ListUser(r *http.Request) *mhttp.ResponseData {
	res := &rpc_impl.User_ListRes{}

	params := &structure.ListUserReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		res.Err = utils.Error_InvalidParams
		return mhttp.NewResponseDataWithError(res)
	}

	client, target, err := getUserClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get user client failed", zap.Error(err))
		res.Err = utils.NewGetClientError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}

	res, err = client.List(context.Background(), (*rpc_impl.User_ListReq)(params))
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_User, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		res.Err = utils.NewGrpcConnectionError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.NewResponseDataWithError(res)
	}

	return mhttp.NewResponseData(res)
}

func CreateUser(r *http.Request) *mhttp.ResponseData {
	res := &rpc_impl.User_CreateRes{}

	params := &structure.CreateUserReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		res.Err = utils.Error_InvalidParams
		return mhttp.NewResponseDataWithError(res)
	}

	client, target, err := getUserClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get user client failed", zap.Error(err))
		res.Err = utils.NewGetClientError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}

	res, err = client.Create(context.Background(), (*rpc_impl.User_CreateReq)(params))
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_User, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		res.Err = utils.NewGrpcConnectionError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.NewResponseDataWithError(res)
	}

	return mhttp.NewResponseData(res)
}

func LockUser(r *http.Request) *mhttp.ResponseData {
	res := &rpc_impl.User_LockRes{}

	params := &structure.LockUserReqParams{}
	params.Decode(r)

	client, target, err := getUserClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get user client failed", zap.Error(err))
		res.Err = utils.NewGrpcConnectionError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}

	res, err = client.Lock(context.Background(), (*rpc_impl.User_LockReq)(params))
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_User, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		res.Err = utils.NewGrpcConnectionError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.NewResponseDataWithError(res)
	}

	return mhttp.NewResponseData(res)
}

func UnlockUser(r *http.Request) *mhttp.ResponseData {
	res := &rpc_impl.User_UnlockRes{}

	params := &structure.UnlockUserReqParams{}
	params.Decode(r)

	client, target, err := getUserClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get user client failed", zap.Error(err))
		res.Err = utils.NewGetClientError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}

	res, err = client.Unlock(context.Background(), (*rpc_impl.User_UnlockReq)(params))
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_User, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		res.Err = utils.NewGrpcConnectionError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.NewResponseDataWithError(res)
	}

	return mhttp.NewResponseData(res)
}

func ModifyUserInfo(r *http.Request) *mhttp.ResponseData {
	res := &rpc_impl.User_ModifyInfoRes{}

	params := &structure.ModifyUserInfoReqParams{}
	params.Decode(r)

	client, target, err := getUserClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get user client failed", zap.Error(err))
		res.Err = utils.NewGetClientError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}

	res, err = client.ModifyInfo(context.Background(), (*rpc_impl.User_ModifyInfoReq)(params))
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_User, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		res.Err = utils.NewGrpcConnectionError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.NewResponseDataWithError(res)
	}

	return mhttp.NewResponseData(res)
}

func ModifyUserPermission(r *http.Request) *mhttp.ResponseData {
	res := &rpc_impl.User_ModifyPermissionRes{}

	params := &structure.ModifyUserPermissionReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		res.Err = utils.Error_InvalidParams
		return mhttp.NewResponseDataWithError(res)
	}

	client, target, err := getUserClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get user client failed", zap.Error(err))
		res.Err = utils.NewGetClientError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}

	res, err = client.ModifyPermission(context.Background(), (*rpc_impl.User_ModifyPermissionReq)(params))
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_User, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		res.Err = utils.NewGrpcConnectionError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.NewResponseDataWithError(res)
	}

	return mhttp.NewResponseData(res)
}
