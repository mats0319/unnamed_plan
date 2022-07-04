package handlers

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/gateway/http/structure_defination"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/http"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/registration_center_embedded"
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
	params := &structure.LoginReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, target, err := getUserClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get user client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.Login(context.Background(), &rpc_impl.User_LoginReq{
		UserName: params.UserName,
		Password: params.Password,
	})
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_User, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(structure.MakeLoginRes(res.UserId, res.Nickname, res.Permission), res.UserId)
}

func ListUser(r *http.Request) *mhttp.ResponseData {
	params := &structure.ListUserReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, target, err := getUserClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get user client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.List(context.Background(), &rpc_impl.User_ListReq{
		OperatorId: params.OperatorID,
		Page: &rpc_impl.Pagination{
			PageSize: uint32(params.PageSize),
			PageNum:  uint32(params.PageNum),
		},
	})
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_User, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(structure.MakeListUserRes(res.Total, usersRPCToHTTP(res.Users...)))
}

func CreateUser(r *http.Request) *mhttp.ResponseData {
	params := &structure.CreateUserReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, target, err := getUserClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get user client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.Create(context.Background(), &rpc_impl.User_CreateReq{
		OperatorId: params.OperatorID,
		UserName:   params.UserName,
		Password:   params.Password,
		Permission: uint32(params.Permission),
	})
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_User, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func LockUser(r *http.Request) *mhttp.ResponseData {
	params := &structure.LockUserReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, target, err := getUserClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get user client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.Lock(context.Background(), &rpc_impl.User_LockReq{
		OperatorId: params.OperatorID,
		UserId:     params.UserID,
	})
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_User, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func UnlockUser(r *http.Request) *mhttp.ResponseData {
	params := &structure.UnlockUserReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, target, err := getUserClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get user client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.Unlock(context.Background(), &rpc_impl.User_UnlockReq{
		OperatorId: params.OperatorID,
		UserId:     params.UserID,
	})
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_User, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func ModifyUserInfo(r *http.Request) *mhttp.ResponseData {
	params := &structure.ModifyUserInfoReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, target, err := getUserClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get user client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.ModifyInfo(context.Background(), &rpc_impl.User_ModifyInfoReq{
		OperatorId: params.OperatorID,
		UserId:     params.UserID,
		CurrPwd:    params.CurrPwd,
		Nickname:   params.Nickname,
		Password:   params.Password,
	})
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_User, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func ModifyUserPermission(r *http.Request) *mhttp.ResponseData {
	params := &structure.ModifyUserPermissionReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, target, err := getUserClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get user client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.ModifyPermission(context.Background(), &rpc_impl.User_ModifyPermissionReq{
		OperatorId: params.OperatorID,
		UserId:     params.UserID,
		Permission: uint32(params.Permission),
	})
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_User, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func usersRPCToHTTP(data ...*rpc_impl.User_Data) []*structure.UserRes {
	res := make([]*structure.UserRes, 0, len(data))
	for i := range data {
		res = append(res, &structure.UserRes{
			UserID:     data[i].UserId,
			UserName:   data[i].UserName,
			Nickname:   data[i].Nickname,
			IsLocked:   data[i].IsLocked,
			Permission: data[i].Permission,
		})
	}

	return res
}
