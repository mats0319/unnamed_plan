package handlers

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/mats9693/unnamed_plan/admin_data/config"
	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/unnamed_plan/admin_data/http/structure_defination"
	"github.com/mats9693/unnamed_plan/admin_data/utils"
	"github.com/mats9693/utils/toy_server/http"
	"github.com/pkg/errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func Login(r *http.Request) string {
	// todo: consider to move into public handler?
	params, err := structure.DecodeLoginParams(r)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	user, err := verifyPwdByUserName(params.Password, params.UserName)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeLoginRes(user.UserID, user.Nickname, user.Permission))
}

func ListUser(r *http.Request) string {
	operatorID := r.PostFormValue("operatorID")
	pageSize, err := strconv.Atoi(r.PostFormValue("pageSize"))
	pageNum, err2 := strconv.Atoi(r.PostFormValue("pageNum"))

	if err != nil || err2 != nil {
		return mhttp.ResponseWithError(errorsToString(err, err2))
	}

	if len(operatorID) < 1 || pageSize < 1 || pageNum < 1 {
		return mhttp.ResponseWithError(error_InvalidParams +
			fmt.Sprintf(", operator id: %s, page size: %d, page num: %d", operatorID, pageSize, pageNum))
	}

	users, count, err := dao.GetUser().QueryPageByPermission(pageSize, pageNum, operatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	userListRes := make([]*structure.UserListRes, 0, len(users))
	for i := range users {
		userListRes = append(userListRes, &structure.UserListRes{
			UserID:     users[i].UserID,
			UserName:   users[i].UserName,
			Nickname:   users[i].Nickname,
			IsLocked:   users[i].IsLocked,
			Permission: users[i].Permission,
			CreatedBy:  users[i].CreatedBy,
		})
	}

	return mhttp.Response(structure.MakeListUserRes(count, userListRes))
}

func CreateUser(r *http.Request) string {
	operatorID := r.PostFormValue("operatorID")
	userName := r.PostFormValue("userName")
	password := r.PostFormValue("password")
	permissionInt, err := strconv.Atoi(r.PostFormValue("permission"))

	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	if len(operatorID) < 1 || len(userName) < 1 || len(password) < 1 {
		return mhttp.ResponseWithError(error_InvalidParams +
			fmt.Sprintf(", operator id: %s, user name: %s, password: %s", operatorID, userName, password))
	}

	permission := uint8(permissionInt)

	operator, err := dao.GetUser().QueryOneInUnlocked(model.User_UserID+" = ?", operatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}
	if operator.Permission < system_config.GetConfiguration().ARankAdminPermission ||
		operator.Permission <= permission {
		return mhttp.ResponseWithError(error_PermissionDenied +
			fmt.Sprintf(", operator: %d, want create: %d", operator.Permission, permission))
	}

	salt := utils.RandomString(10)
	err = dao.GetUser().Insert(&model.User{
		UserName:   userName,
		Nickname:   userName,
		Password:   utils.CalcSHA256(password, salt),
		Salt:       salt,
		Permission: permission,
		CreatedBy:  operatorID,
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeCreateUserRes(true))
}

func LockUser(r *http.Request) string {
	operatorID := r.PostFormValue("operatorID")
	userID := r.PostFormValue("userID")

	if len(operatorID) < 1 || len(userID) < 1 || operatorID == userID {
		return mhttp.ResponseWithError(error_InvalidParams +
			fmt.Sprintf(", operator id: %s, user id: %s", operatorID, userID))
	}

	users, err := dao.GetUser().Query(model.User_UserID+" in (?)", pg.In([]string{operatorID, userID}))
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	users, err = sortUsersByUserID(users, []string{operatorID, userID})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	if users[1].IsLocked {
		return mhttp.ResponseWithError(error_UserLocked)
	}
	if users[0].Permission <= users[1].Permission ||
		users[0].Permission < system_config.GetConfiguration().ARankAdminPermission {
		return mhttp.ResponseWithError(error_PermissionDenied +
			fmt.Sprintf(", operator: %d, user: %d, need: %d",
				users[0].Permission, users[1].Permission, system_config.GetConfiguration().ARankAdminPermission))
	}

	users[1].IsLocked = true

	err = dao.GetUser().UpdateColumnsByUserID(users[1], model.User_IsLocked)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeLockUserRes(true))
}

func UnlockUser(r *http.Request) string {
	operatorID := r.PostFormValue("operatorID")
	userID := r.PostFormValue("userID")

	if len(operatorID) < 1 || len(userID) < 1 || operatorID == userID {
		return mhttp.ResponseWithError(error_InvalidParams +
			fmt.Sprintf(", operator id: %s, user id: %s", operatorID, userID))
	}

	users, err := dao.GetUser().Query(model.User_UserID+" in (?)", pg.In([]string{operatorID, userID}))
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	users, err = sortUsersByUserID(users, []string{operatorID, userID})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	if !users[1].IsLocked {
		return mhttp.ResponseWithError(error_UserUnlocked)
	}
	if users[0].Permission <= users[1].Permission ||
		users[0].Permission < system_config.GetConfiguration().ARankAdminPermission {
		return mhttp.ResponseWithError(error_PermissionDenied +
			fmt.Sprintf(", operator: %d, user: %d, need: %d",
				users[0].Permission, users[1].Permission, system_config.GetConfiguration().ARankAdminPermission))
	}

	users[1].IsLocked = false

	err = dao.GetUser().UpdateColumnsByUserID(users[1], model.User_IsLocked)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeUnlockUserRes(true))
}

func ModifyUserInfo(r *http.Request) string {
	operatorID := r.PostFormValue("operatorID")
	userID := r.PostFormValue("userID")
	currPwd := r.PostFormValue("currPwd")
	nickname := r.PostFormValue("nickname")
	password := r.PostFormValue("password")

	if len(operatorID) < 1 || len(userID) < 1 || operatorID != userID {
		return mhttp.ResponseWithError(error_InvalidParams +
			fmt.Sprintf(", operator id: %s, user id: %s", operatorID, userID))
	}

	if len(nickname)+len(password) < 1 {
		return mhttp.ResponseWithError(error_NoValidModification)
	}

	user, err := verifyPwdByUserID(currPwd, userID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	updateColumns := make([]string, 0, 2)
	if len(nickname) > 0 {
		user.Nickname = nickname
		updateColumns = append(updateColumns, model.User_Nickname)
	}
	if len(password) > 0 {
		user.Password = utils.CalcSHA256(password, user.Salt)
		updateColumns = append(updateColumns, model.User_Password)
	}

	err = dao.GetUser().UpdateColumnsByUserID(user, updateColumns...)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeModifyUserInfoRes(true))
}

func ModifyUserPermission(r *http.Request) string {
	operatorID := r.PostFormValue("operatorID")
	userID := r.PostFormValue("userID")
	permissionInt, err := strconv.Atoi(r.PostFormValue("permission"))

	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	if len(operatorID) < 1 || len(userID) < 1 || operatorID == userID {
		return mhttp.ResponseWithError(error_InvalidParams +
			fmt.Sprintf(", operator id: %s, user id: %s", operatorID, userID))
	}

	users, err := dao.GetUser().Query(model.User_UserID+" in (?)", pg.In([]string{operatorID, userID}))
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	users, err = sortUsersByUserID(users, []string{operatorID, userID})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	permission := uint8(permissionInt)

	if users[0].Permission < system_config.GetConfiguration().SRankAdminPermission ||
		users[1].Permission >= system_config.GetConfiguration().SRankAdminPermission ||
		permission >= system_config.GetConfiguration().SRankAdminPermission {
		return mhttp.ResponseWithError(error_PermissionDenied +
			fmt.Sprintf(", operator: %d, user: %d, user new: %d", users[0].Permission, users[1].Permission, permission))
	}

	users[1].Permission = permission

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
