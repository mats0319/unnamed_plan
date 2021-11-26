package handlers

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/gateway/http/structure_defination"
	"github.com/mats9693/unnamed_plan/services/gateway/rpc"
	"github.com/mats9693/unnamed_plan/shared/proto/impl"
	mconst "github.com/mats9693/utils/toy_server/const"
	"github.com/mats9693/utils/toy_server/http"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func Login(r *http.Request) *mhttp.ResponseData {
	params := &structure.LoginReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	res, err := rpc.GetRPCClient().UserClient.Login(context.Background(), &rpc_impl.User_LoginReq{
		UserName: params.UserName,
		Password: params.Password,
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeLoginRes(res.UserId, res.Nickname, res.Permission), res.UserId)
}

func ListUser(r *http.Request) *mhttp.ResponseData {
	params := &structure.ListUserReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	res, err := rpc.GetRPCClient().UserClient.List(context.Background(), &rpc_impl.User_ListReq{
		OperatorId: params.OperatorID,
		Page:       &rpc_impl.Pagination{
			PageSize: uint32(params.PageSize),
			PageNum:  uint32(params.PageNum),
		},
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeListUserRes(res.Total, usersRPCToHTTP(res.Users...)))
}

func CreateUser(r *http.Request) *mhttp.ResponseData {
	params := &structure.CreateUserReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	_, err := rpc.GetRPCClient().UserClient.Create(context.Background(), &rpc_impl.User_CreateReq{
		OperatorId: params.OperatorID,
		UserName:   params.UserName,
		Password:   params.Password,
		Permission: uint32(params.Permission),
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func LockUser(r *http.Request) *mhttp.ResponseData {
	params := &structure.LockUserReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	_, err := rpc.GetRPCClient().UserClient.Lock(context.Background(), &rpc_impl.User_LockReq{
		OperatorId: params.OperatorID,
		UserId:     params.UserID,
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func UnlockUser(r *http.Request) *mhttp.ResponseData {
	params := &structure.UnlockUserReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	_, err := rpc.GetRPCClient().UserClient.Unlock(context.Background(), &rpc_impl.User_UnlockReq{
		OperatorId: params.OperatorID,
		UserId:     params.UserID,
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func ModifyUserInfo(r *http.Request) *mhttp.ResponseData {
	params := &structure.ModifyUserInfoReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	_, err := rpc.GetRPCClient().UserClient.ModifyInfo(context.Background(), &rpc_impl.User_ModifyInfoReq{
		OperatorId: params.OperatorID,
		UserId:     params.UserID,
		CurrPwd:    params.CurrPwd,
		Nickname:   params.Nickname,
		Password:   params.Password,
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func ModifyUserPermission(r *http.Request) *mhttp.ResponseData {
	params := &structure.ModifyUserPermissionReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	_, err := rpc.GetRPCClient().UserClient.ModifyPermission(context.Background(), &rpc_impl.User_ModifyPermissionReq{
		OperatorId: params.OperatorID,
		UserId:     params.UserID,
		Permission: uint32(params.Permission),
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
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
			CreatedBy:  data[i].CreatedBy,
		})
	}

	return res
}