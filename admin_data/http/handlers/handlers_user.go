package handlers

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/mats9693/unnamed_plan/admin_data/config"
	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/unnamed_plan/admin_data/http/response_type"
	"github.com/mats9693/unnamed_plan/admin_data/kits"
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

func Login(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	userName := r.PostFormValue("userName")
	password := r.PostFormValue("password")

	user, err := checkPwdByUserName(password, userName)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	resData := &struct {
		UserID     string `json:"userID"`
		Nickname   string `json:"nickname"`
		Permission uint8  `json:"permission"`
	}{
		UserID:     user.UserID,
		Nickname:   user.Nickname,
		Permission: user.Permission,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}

func ListUser(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	pageSize, err := strconv.Atoi(r.PostFormValue("pageSize"))
	pageNum, err2 := strconv.Atoi(r.PostFormValue("pageNum"))

	if err != nil || err2 != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()+err2.Error()))
		return
	}

	if len(operatorID) < 1 || pageSize < 1 || pageNum < 1 {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid param, operator id: %s, page size: %d, page num: %d.", operatorID, pageSize, pageNum)))
		return
	}

	users, count, err := dao.GetUser().QueryPageByPermission(pageSize, pageNum, operatorID)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	usersRes := make([]*http_res_type.HttpResUser, 0, len(users))
	for i := range users {
		usersRes = append(usersRes, &http_res_type.HttpResUser{
			UserID:     users[i].UserID,
			UserName:   users[i].UserName,
			Nickname:   users[i].Nickname,
			IsLocked:   users[i].IsLocked,
			Permission: users[i].Permission,
			CreatedBy:  users[i].CreatedBy,
		})
	}

	resData := &struct {
		Total int                          `json:"total"`
		Users []*http_res_type.HttpResUser `json:"users"`
	}{
		Total: count,
		Users: usersRes,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}

func ModifyUserInfo(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	userID := r.PostFormValue("userID")
	currPwd := r.PostFormValue("currPwd")
	nickname := r.PostFormValue("nickname")
	password := r.PostFormValue("password")

	if len(operatorID) < 1 || len(userID) < 1 || operatorID != userID {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, user id: %s", operatorID, userID)))
		return
	}

	if len(nickname)+len(password) < 1 {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError("invalid params, not any modification received"))
		return
	}

	user, err := checkPwdByUserID(currPwd, userID)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	updateColumns := make([]string, 0, 2)
	if len(nickname) > 0 {
		user.Nickname = nickname
		updateColumns = append(updateColumns, model.User_Nickname)
	}
	if len(password) > 0 {
		user.Password = kits.CalcSHA256(password, user.Salt)
		updateColumns = append(updateColumns, model.User_Password)
	}

	err = dao.GetUser().UpdateColumnsByUserID(user, updateColumns...)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	resData := &struct {
		IsSuccess bool `json:"isSuccess"`
	}{
		IsSuccess: true,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	userName := r.PostFormValue("userName")
	password := r.PostFormValue("password")
	permissionInt, err := strconv.Atoi(r.PostFormValue("permission"))

	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	if len(operatorID) < 1 || len(userName) < 1 || len(password) < 1 {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, user name: %s, password: %s", operatorID, userName, password)))
		return
	}

	permission := uint8(permissionInt)

	operator, err := dao.GetUser().QueryUnlocked(model.User_UserID+" = ?", operatorID)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}
	if operator.Permission < system_config.GetConfiguration().ARankAdminPermission ||
		operator.Permission <= permission {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("permission denied, operator: %d, want create: %d.",
			operator.Permission, permission)))
		return
	}

	salt := kits.RandomString(10)
	err = dao.GetUser().Insert(&model.User{
		UserName:   userName,
		Nickname:   userName,
		Password:   kits.CalcSHA256(password, salt),
		Salt:       salt,
		Permission: permission,
		CreatedBy:  operatorID,
	})
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	resData := &struct {
		IsSuccess bool `json:"isSuccess"`
	}{
		IsSuccess: true,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}

func LockUser(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	userID := r.PostFormValue("userID")

	if len(operatorID) < 1 || len(userID) < 1 || operatorID == userID {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, user id: %s", operatorID, userID)))
		return
	}

	users, err := dao.GetUser().Query(model.User_UserID+" in (?)", pg.In([]string{operatorID, userID}))
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	users, err = sortUsersByUserID(users, []string{operatorID, userID})
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	if users[1].IsLocked {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError("user is already locked"))
		return
	}
	if users[0].Permission <= users[1].Permission ||
		users[0].Permission < system_config.GetConfiguration().ARankAdminPermission {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("permission denied, operator: %d, user: %d, need: %d",
			users[0].Permission, users[1].Permission, system_config.GetConfiguration().SRankAdminPermission)))
		return
	}

	users[1].IsLocked = true

	err = dao.GetUser().UpdateColumnsByUserID(users[1], model.User_IsLocked)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	resData := &struct {
		IsSuccess bool `json:"isSuccess"`
	}{
		IsSuccess: true,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}

func UnlockUser(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	userID := r.PostFormValue("userID")

	if len(operatorID) < 1 || len(userID) < 1 || operatorID == userID {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, user id: %s", operatorID, userID)))
		return
	}

	users, err := dao.GetUser().Query(model.User_UserID+" in (?)", pg.In([]string{operatorID, userID}))
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	users, err = sortUsersByUserID(users, []string{operatorID, userID})
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	if !users[1].IsLocked {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError("user is already unlocked"))
		return
	}
	if users[0].Permission <= users[1].Permission ||
		users[0].Permission < system_config.GetConfiguration().ARankAdminPermission {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("permission denied, operator: %d, user: %d, need: %d",
			users[0].Permission, users[1].Permission, system_config.GetConfiguration().SRankAdminPermission)))
		return
	}

	users[1].IsLocked = false

	err = dao.GetUser().UpdateColumnsByUserID(users[1], model.User_IsLocked)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	resData := &struct {
		IsSuccess bool `json:"isSuccess"`
	}{
		IsSuccess: true,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}

func ModifyUserPermission(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	userID := r.PostFormValue("userID")
	permissionInt, err := strconv.Atoi(r.PostFormValue("permission"))

	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	if len(operatorID) < 1 || len(userID) < 1 || operatorID == userID {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, user id: %s", operatorID, userID)))
		return
	}

	users, err := dao.GetUser().Query(model.User_UserID+" in (?)", pg.In([]string{operatorID, userID}))
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	users, err = sortUsersByUserID(users, []string{operatorID, userID})
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	permission := uint8(permissionInt)

	if users[0].Permission < system_config.GetConfiguration().SRankAdminPermission ||
		users[1].Permission >= system_config.GetConfiguration().SRankAdminPermission ||
		permission >= system_config.GetConfiguration().SRankAdminPermission {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("permission denied, operator: %d, user: %d, user new: %d",
			users[0].Permission, users[1].Permission, permission)))
		return
	}

	users[1].Permission = permission

	err = dao.GetUser().UpdateColumnsByUserID(users[1], model.User_Permission)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	resData := &struct {
		IsSuccess bool `json:"isSuccess"`
	}{
		IsSuccess: true,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
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
