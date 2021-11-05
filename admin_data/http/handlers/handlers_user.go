package handlers

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/mats9693/unnamed_plan/admin_data/config"
	"github.com/mats9693/unnamed_plan/admin_data/http/structure_defination"
	"github.com/mats9693/unnamed_plan/admin_data/utils"
	"github.com/mats9693/unnamed_plan/shared/db/dao"
	"github.com/mats9693/unnamed_plan/shared/db/model"
	"github.com/mats9693/utils/toy_server/http"
	"github.com/mats9693/utils/toy_server/utils"
	"github.com/pkg/errors"
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

	user, err := verifyPwdByUserName(params.Password, params.UserName)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeLoginRes(user.UserID, user.Nickname, user.Permission), user.UserID)
}

func ListUser(r *http.Request) *mhttp.ResponseData {
	params := &structure.ListUserReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	if len(params.OperatorID) < 1 || params.PageSize < 1 || params.PageNum < 1 {
		return mhttp.ResponseWithError(error_InvalidParams,
			String("operator id", params.OperatorID),
			Int("page size", params.PageSize),
			Int("page num", params.PageNum))
	}

	users, count, err := dao.GetUser().QueryPageByPermission(params.PageSize, params.PageNum, params.OperatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	userListRes := make([]*structure.UserRes, 0, len(users))
	for i := range users {
		userListRes = append(userListRes, userDBToHTTPRes(users[i]))
	}

	return mhttp.Response(structure.MakeListUserRes(count, userListRes))
}

func CreateUser(r *http.Request) *mhttp.ResponseData {
	params := &structure.CreateUserReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	if len(params.OperatorID) < 1 || len(params.UserName) < 1 || len(params.Password) < 1 {
		return mhttp.ResponseWithError(error_InvalidParams,
			String("operator id", params.OperatorID),
			String("user name", params.UserName),
			String("password", params.Password))
	}

	operator, err := dao.GetUser().QueryOneInUnlocked(model.User_UserID+" = ?", params.OperatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}
	if operator.Permission < system_config.GetConfiguration().ARankAdminPermission ||
		operator.Permission <= params.Permission {
		return mhttp.ResponseWithError(error_PermissionDenied,
			Uint8("operator", operator.Permission),
			Uint8("want create", params.Permission))
	}

	salt := mutils.RandomHexString(10)
	err = dao.GetUser().Insert(&model.User{
		UserName:   params.UserName,
		Nickname:   params.UserName,
		Password:   utils.CalcSHA256(params.Password, salt),
		Salt:       salt,
		Permission: params.Permission,
		CreatedBy:  params.OperatorID,
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeCreateUserRes(true))
}

func LockUser(r *http.Request) *mhttp.ResponseData {
	params := &structure.LockUserReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	if len(params.OperatorID) < 1 || len(params.UserID) < 1 || params.OperatorID == params.UserID {
		return mhttp.ResponseWithError(error_InvalidParams,
			String("operator id", params.OperatorID),
			String("user id", params.UserID))
	}

	users, err := dao.GetUser().Query(model.User_UserID+" in (?)", pg.In([]string{params.OperatorID, params.UserID}))
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	users, err = sortUsersByUserID(users, []string{params.OperatorID, params.UserID})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	if users[1].IsLocked {
		return mhttp.ResponseWithError(error_UserLocked)
	}
	if users[0].Permission <= users[1].Permission ||
		users[0].Permission < system_config.GetConfiguration().ARankAdminPermission {
		return mhttp.ResponseWithError(error_PermissionDenied,
			Uint8("operator", users[0].Permission),
			Uint8("user", users[1].Permission),
			Uint8("need", system_config.GetConfiguration().ARankAdminPermission))
	}

	users[1].IsLocked = true

	err = dao.GetUser().UpdateColumnsByUserID(users[1], model.User_IsLocked)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeLockUserRes(true))
}

func UnlockUser(r *http.Request) *mhttp.ResponseData {
	params := &structure.UnlockUserReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	if len(params.OperatorID) < 1 || len(params.UserID) < 1 || params.OperatorID == params.UserID {
		return mhttp.ResponseWithError(error_InvalidParams,
			String("operator id", params.OperatorID),
			String("user id", params.UserID))
	}

	users, err := dao.GetUser().Query(model.User_UserID+" in (?)", pg.In([]string{params.OperatorID, params.UserID}))
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	users, err = sortUsersByUserID(users, []string{params.OperatorID, params.UserID})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	if !users[1].IsLocked {
		return mhttp.ResponseWithError(error_UserUnlocked)
	}
	if users[0].Permission <= users[1].Permission ||
		users[0].Permission < system_config.GetConfiguration().ARankAdminPermission {
		return mhttp.ResponseWithError(error_PermissionDenied,
			Uint8("operator", users[0].Permission),
			Uint8("user", users[1].Permission),
			Uint8("need", system_config.GetConfiguration().ARankAdminPermission))
	}

	users[1].IsLocked = false

	err = dao.GetUser().UpdateColumnsByUserID(users[1], model.User_IsLocked)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeUnlockUserRes(true))
}

func ModifyUserInfo(r *http.Request) *mhttp.ResponseData {
	params := &structure.ModifyUserInfoReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	if len(params.OperatorID) < 1 || len(params.UserID) < 1 || params.OperatorID != params.UserID {
		return mhttp.ResponseWithError(error_InvalidParams,
			String("operator id", params.OperatorID),
			String("user id", params.UserID))
	}

	if len(params.Nickname)+len(params.Password) < 1 {
		return mhttp.ResponseWithError(error_NoValidModification)
	}

	user, err := verifyPwdByUserID(params.CurrPwd, params.UserID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	updateColumns := make([]string, 0, 2)
	if len(params.Nickname) > 0 {
		user.Nickname = params.Nickname
		updateColumns = append(updateColumns, model.User_Nickname)
	}
	if len(params.Password) > 0 {
		user.Password = utils.CalcSHA256(params.Password, user.Salt)
		updateColumns = append(updateColumns, model.User_Password)
	}

	err = dao.GetUser().UpdateColumnsByUserID(user, updateColumns...)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeModifyUserInfoRes(true))
}

func ModifyUserPermission(r *http.Request) *mhttp.ResponseData {
	params := &structure.ModifyUserPermissionReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	if len(params.OperatorID) < 1 || len(params.UserID) < 1 || params.OperatorID == params.UserID {
		return mhttp.ResponseWithError(error_InvalidParams,
			String("operator id", params.OperatorID),
			String("user id", params.UserID))
	}

	users, err := dao.GetUser().Query(model.User_UserID+" in (?)", pg.In([]string{params.OperatorID, params.UserID}))
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	users, err = sortUsersByUserID(users, []string{params.OperatorID, params.UserID})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	if users[0].Permission < system_config.GetConfiguration().SRankAdminPermission ||
		users[1].Permission >= system_config.GetConfiguration().SRankAdminPermission ||
		params.Permission >= system_config.GetConfiguration().SRankAdminPermission {
		return mhttp.ResponseWithError(error_PermissionDenied,
			Uint8("operator", users[0].Permission),
			Uint8("user", users[1].Permission),
			Uint8("user new", params.Permission))
	}

	users[1].Permission = params.Permission

	err = dao.GetUser().UpdateColumnsByUserID(users[1], model.User_Permission)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeModifyUserPermissionRes(true))
}

func sortUsersByUserID(users []*model.User, order []string) ([]*model.User, error) {
	if len(users) != len(order) {
		return nil, errors.New(fmt.Sprintf("unmatched users amount, users %d, orders %d", len(users), len(order)))
	}

	length := len(users)
	for i := 0; i < length; i++ {
		for j := i; j < length; j++ {
			if order[j] == users[i].UserID {
				users[i], users[j] = users[j], users[i]
				break
			}
		}
	}

	unmatchedIndex := -1
	for i := 0; i < length; i++ {
		if users[i].UserID != order[i] {
			unmatchedIndex = i
			break
		}
	}

	if unmatchedIndex >= 0 {
		return nil, errors.New(fmt.Sprintf("unmatched user id: %s", users[unmatchedIndex].UserID))
	}

	return users, nil
}

func userDBToHTTPRes(data *model.User) *structure.UserRes {
	if data == nil {
		return &structure.UserRes{}
	}

	return &structure.UserRes{
		UserID:     data.UserID,
		UserName:   data.UserName,
		Nickname:   data.Nickname,
		IsLocked:   data.IsLocked,
		Permission: data.Permission,
		CreatedBy:  data.CreatedBy,
	}
}
