package rpc

import (
	"context"
	"fmt"
	"github.com/mats9693/unnamed_plan/services/1_user/config"
	"github.com/mats9693/unnamed_plan/services/1_user/db"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type userServerImpl struct {
	rpc_impl.UnimplementedIUserServer
}

var userServerImplIns = &userServerImpl{}

func GetUserServer() rpc_impl.IUserServer {
	return userServerImplIns
}

func (s *userServerImpl) Login(_ context.Context, req *rpc_impl.User_LoginReq) (*rpc_impl.User_LoginRes, error) {
	res := &rpc_impl.User_LoginRes{}

	if len(req.UserName) < 1 || len(req.Password) < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("user name", req.UserName),
			zap.String("password", req.Password))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	user, err := db.GetUserDao().QueryOneByUserName(req.UserName)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	if user.Password != utils.CalcSHA256(req.Password, user.Salt) {
		mlog.Logger().Error(mconst.Error_InvalidAccountOrPassword)
		res.Err = utils.Error_InvalidAccountOrPassword.ToRPC()
		return res, nil
	}

	res.UserId = user.ID
	res.Nickname = user.Nickname
	res.Permission = uint32(user.Permission)

	return res, nil
}

func (s *userServerImpl) List(_ context.Context, req *rpc_impl.User_ListReq) (*rpc_impl.User_ListRes, error) {
	res := &rpc_impl.User_ListRes{}

	if len(req.OperatorId) < 1 || req.GetPage() == nil || req.Page.PageSize < 1 || req.Page.PageNum < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("operator", req.OperatorId),
			zap.String("page info", req.Page.String()))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	pageSize := int(req.Page.PageSize)
	pageNum := int(req.Page.PageNum)

	users, count, err := db.GetUserDao().QueryPageLEPermission(pageSize, pageNum, req.OperatorId)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	res.Total = uint32(count)
	res.Users = usersDBToRPC(users...)

	return res, nil
}

func (s *userServerImpl) Create(_ context.Context, req *rpc_impl.User_CreateReq) (*rpc_impl.User_CreateRes, error) {
	res := &rpc_impl.User_CreateRes{}

	if len(req.OperatorId) < 1 || len(req.UserName) < 1 || len(req.Password) < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("operator", req.OperatorId),
			zap.String("user name", req.UserName),
			zap.String("password", req.Password))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	operator, err := db.GetUserDao().QueryOne(req.OperatorId)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	permission := uint8(req.Permission)
	if operator.Permission < config.GetConfig().ARankAdminPermission || operator.Permission <= permission {
		mlog.Logger().Error(mconst.Error_PermissionDenied)
		res.Err = utils.Error_PermissionDenied.ToRPC()
		return res, nil
	}

	salt := utils.RandomHexString(10)
	err = db.GetUserDao().Insert(&model.User{
		UserName:   req.UserName,
		Nickname:   req.UserName,
		Password:   utils.CalcSHA256(req.Password, salt),
		Salt:       salt,
		Permission: permission,
		CreatedBy:  req.OperatorId,
	})
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	return res, nil
}

func (s *userServerImpl) Lock(_ context.Context, req *rpc_impl.User_LockReq) (*rpc_impl.User_LockRes, error) {
	res := &rpc_impl.User_LockRes{}

	if len(req.OperatorId) < 1 || len(req.UserId) < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("operator", req.OperatorId),
			zap.String("user", req.UserId))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	users, err := db.GetUserDao().Query([]string{req.OperatorId, req.UserId})
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	users, err = sortUsersByUserID(users, []string{req.OperatorId, req.UserId})
	if err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.Error(err))
		res.Err = utils.NewExecError(err.Error()).ToRPC()
		return res, nil
	}

	if users[1].IsLocked {
		mlog.Logger().Error(mconst.Error_UserAlreadyLocked)
		res.Err = utils.NewExecError(mconst.Error_UserAlreadyLocked).ToRPC()
		return res, nil
	}
	if users[0].Permission < config.GetConfig().ARankAdminPermission || users[0].Permission <= users[1].Permission {
		mlog.Logger().Error(mconst.Error_PermissionDenied)
		res.Err = utils.NewExecError(mconst.Error_PermissionDenied).ToRPC()
		return res, nil
	}

	users[1].IsLocked = true

	err = db.GetUserDao().UpdateColumnsByUserID(users[1], model.User_IsLocked)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	return res, nil
}

func (s *userServerImpl) Unlock(_ context.Context, req *rpc_impl.User_UnlockReq) (*rpc_impl.User_UnlockRes, error) {
	res := &rpc_impl.User_UnlockRes{}

	if len(req.OperatorId) < 1 || len(req.UserId) < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("operator", req.OperatorId),
			zap.String("user", req.UserId))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	users, err := db.GetUserDao().Query([]string{req.OperatorId, req.UserId})
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	users, err = sortUsersByUserID(users, []string{req.OperatorId, req.UserId})
	if err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.Error(err))
		res.Err = utils.NewExecError(err.Error()).ToRPC()
		return res, nil
	}

	if !users[1].IsLocked {
		mlog.Logger().Error(mconst.Error_UserAlreadyUnlocked)
		res.Err = utils.NewExecError(mconst.Error_UserAlreadyUnlocked).ToRPC()
		return res, nil
	}
	if users[0].Permission < config.GetConfig().ARankAdminPermission || users[0].Permission <= users[1].Permission {
		mlog.Logger().Error(mconst.Error_PermissionDenied)
		res.Err = utils.NewExecError(mconst.Error_PermissionDenied).ToRPC()
		return res, nil
	}

	users[1].IsLocked = false

	err = db.GetUserDao().UpdateColumnsByUserID(users[1], model.User_IsLocked)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	return res, nil
}

func (s *userServerImpl) ModifyInfo(_ context.Context, req *rpc_impl.User_ModifyInfoReq) (*rpc_impl.User_ModifyInfoRes, error) {
	res := &rpc_impl.User_ModifyInfoRes{}

	if len(req.OperatorId) < 1 || len(req.UserId) < 1 || req.OperatorId != req.UserId {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("operator", req.OperatorId),
			zap.String("user", req.UserId))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	if len(req.Nickname)+len(req.Password) < 1 {
		mlog.Logger().Error(mconst.Error_NoValidModification)
		res.Err = utils.Error_NoValidModification.ToRPC()
		return res, nil
	}

	user, err := verifyPwdByUserID(req.UserId, req.CurrPwd)
	if err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError)
		res.Err = utils.NewExecError(err.Error()).ToRPC()
		return res, nil
	}

	updateColumns := make([]string, 0, 2)
	if len(req.Nickname) > 0 {
		user.Nickname = req.Nickname
		updateColumns = append(updateColumns, model.User_Nickname)
	}
	if len(req.Password) > 0 {
		user.Password = utils.CalcSHA256(req.Password, user.Salt)
		updateColumns = append(updateColumns, model.User_Password)
	}

	err = db.GetUserDao().UpdateColumnsByUserID(user, updateColumns...)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	return res, nil
}

func (s *userServerImpl) ModifyPermission(_ context.Context, req *rpc_impl.User_ModifyPermissionReq) (*rpc_impl.User_ModifyPermissionRes, error) {
	res := &rpc_impl.User_ModifyPermissionRes{}

	if len(req.OperatorId) < 1 || len(req.UserId) < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("operator", req.OperatorId),
			zap.String("user", req.UserId))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	users, err := db.GetUserDao().Query([]string{req.OperatorId, req.UserId})
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	users, err = sortUsersByUserID(users, []string{req.OperatorId, req.UserId})
	if err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.Error(err))
		res.Err = utils.NewExecError(err.Error()).ToRPC()
		return res, nil
	}

	permission := uint8(req.Permission)
	if users[0].Permission < config.GetConfig().SRankAdminPermission ||
		users[1].Permission >= config.GetConfig().SRankAdminPermission ||
		permission >= config.GetConfig().SRankAdminPermission {
		mlog.Logger().Error(mconst.Error_PermissionDenied)
		res.Err = utils.NewExecError(mconst.Error_PermissionDenied).ToRPC()
		return res, nil
	}

	users[1].Permission = permission

	err = db.GetUserDao().UpdateColumnsByUserID(users[1], model.User_Permission)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	return res, nil
}

func (s *userServerImpl) Authenticate(_ context.Context, req *rpc_impl.User_AuthenticateReq) (*rpc_impl.User_AuthenticateRes, error) {
	res := &rpc_impl.User_AuthenticateRes{}

	if len(req.UserId) < 1 || len(req.Password) < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("user", req.UserId),
			zap.String("password", req.Password))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	_, err := verifyPwdByUserID(req.UserId, req.Password)
	if err != nil {
		mlog.Logger().Error(mconst.Error_InvalidAccountOrPassword)
		res.Err = utils.Error_InvalidAccountOrPassword.ToRPC()
		return res, nil
	}

	return res, nil
}

func usersDBToRPC(data ...*model.User) []*rpc_impl.User_Data {
	res := make([]*rpc_impl.User_Data, 0, len(data))
	for i := range data {
		res = append(res, &rpc_impl.User_Data{
			UserId:     data[i].ID,
			UserName:   data[i].UserName,
			Nickname:   data[i].Nickname,
			IsLocked:   data[i].IsLocked,
			Permission: uint32(data[i].Permission),
		})
	}

	return res
}

func verifyPwdByUserID(userID string, password string) (*model.User, error) {
	user, err := db.GetUserDao().QueryOne(userID)
	if err != nil {
		return nil, err
	}

	if user.Password != utils.CalcSHA256(password, user.Salt) {
		return nil, utils.NewError(mconst.Error_InvalidAccountOrPassword)
	}

	return user, nil
}

func sortUsersByUserID(users []*model.User, order []string) ([]*model.User, error) {
	if len(users) != len(order) {
		return nil, errors.New(fmt.Sprintf("unmatched users amount, users %d, orders %d", len(users), len(order)))
	}

	length := len(users)
	for i := 0; i < length; i++ {
		for j := i; j < length; j++ {
			if order[j] == users[i].ID {
				users[i], users[j] = users[j], users[i]
				break
			}
		}
	}

	unmatchedIndex := -1
	for i := 0; i < length; i++ {
		if users[i].ID != order[i] {
			unmatchedIndex = i
			break
		}
	}

	if unmatchedIndex >= 0 {
		return nil, errors.New(fmt.Sprintf("unmatched user id: %s", users[unmatchedIndex].ID))
	}

	return users, nil
}
