package http

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/mats9693/unnamed_plan/admin_data/config"
	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/unnamed_plan/admin_data/kits"
	"github.com/mats9693/utils/toy_server/http"
	"net/http"
	"strconv"
)

func login(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	userName := r.PostFormValue("userName")
	password := r.PostFormValue("password")

	user, err := dao.GetUser().QueryOne(model.User_UserName+" = ?", userName)
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	if user.Password != kits.CalcPassword(password, user.Salt) {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError("invalid account or password"))
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

	_, _ = fmt.Fprintln(w, shttp.Response(resData))

	return
}

func listUser(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	userID := r.PostFormValue("operatorID")
	pageSize, pageSizeErr := strconv.Atoi(r.PostFormValue("pageSize"))
	pageNum, pageNumErr := strconv.Atoi(r.PostFormValue("pageNum"))

	if pageSizeErr != nil || pageNumErr != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(pageSizeErr.Error()+pageNumErr.Error()))
		return
	}

	if len(userID) < 1 || pageSize < 1 || pageNum < 1 {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("invalid param, userID: %s, page size: %d, page num: %d.", userID, pageSize, pageNum)))
		return
	}

	users, count, err := dao.GetUser().QueryPageByPermission(pageSize, pageNum, userID)
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	type HttpResUser struct {
		UserID     string `json:"userID"`
		Nickname   string `json:"nickname"`
		IsLocked   bool   `json:"isLocked"`
		Permission uint8  `json:"permission"`
	}

	var usersRes []*HttpResUser
	for i := range users {
		if users[i].UserID != userID {
			usersRes = append(usersRes, &HttpResUser{
				UserID:     users[i].UserID,
				Nickname:   users[i].Nickname,
				IsLocked:   users[i].IsLocked,
				Permission: users[i].Permission,
			})
		}
	}

	resData := &struct {
		Total int            `json:"total"`
		Users []*HttpResUser `json:"users"`
	}{
		Total: count,
		Users: usersRes,
	}

	_, _ = fmt.Fprintln(w, shttp.Response(resData))

	return
}

func modifyUserInfo(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	userID := r.PostFormValue("userID")
	currPwd := r.PostFormValue("currPwd")
	nickname := r.PostFormValue("nickname")
	password := r.PostFormValue("password")

	if operatorID != userID {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError("invalid params, operator id is not equal to user id"))
		return
	}

	if len(nickname)+len(password) < 1 {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError("invalid params, not any modification received"))
		return
	}

	user, err := dao.GetUser().QueryOne(model.User_UserID+" = ?", userID)
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	} else if user.Password != kits.CalcPassword(currPwd, user.Salt) {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError("invalid password"))
		return
	}

	updateColumns := make([]string, 0, 2)
	if len(nickname) > 0 {
		user.Nickname = nickname
		updateColumns = append(updateColumns, model.User_Nickname)
	}
	if len(password) > 0 {
		user.Password = kits.CalcPassword(password, user.Salt)
		updateColumns = append(updateColumns, model.User_Password)
	}

	err = dao.GetUser().UpdateColumnsByUserID(user, updateColumns...)
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	resData := &struct {
		IsSuccess bool `json:"isSuccess"`
	}{
		IsSuccess: true,
	}

	_, _ = fmt.Fprintln(w, shttp.Response(resData))

	return
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	userName := r.PostFormValue("userName")
	password := r.PostFormValue("password")
	permissionInt, err := strconv.Atoi(r.PostFormValue("permission"))

	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	if len(operatorID) < 1 || len(userName) < 1 || len(password) < 1 {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, user name: %s, password: %s", operatorID, userName, password)))
		return
	}

	permission := uint8(permissionInt)

	operator, err := dao.GetUser().QueryOne(model.User_UserID+" = ?", operatorID)
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	} else if operator.Permission < system_config.GetConfiguration().ARankAdminPermission ||
		operator.Permission <= permission {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("permission denied, operator: %d, want create: %d.",
			operator.Permission, permission)))
		return
	}

	salt := kits.RandomString(10)
	err = dao.GetUser().Insert(&model.User{
		UserName:   userName,
		Nickname:   userName,
		Password:   kits.CalcPassword(password, salt),
		Salt:       salt,
		Permission: permission,
		CreatedBy:  operatorID,
	})
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	resData := &struct {
		IsSuccess bool `json:"isSuccess"`
	}{
		IsSuccess: true,
	}

	_, _ = fmt.Fprintln(w, shttp.Response(resData))

	return
}

func lockUser(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	userID := r.PostFormValue("userID")

	if len(operatorID) < 1 || len(userID) < 1 || operatorID == userID {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, user id: %s", operatorID, userID)))
		return
	}

	users, err := dao.GetUser().Query(model.User_UserID+" in (?)", pg.In([]string{operatorID, userID}))
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	} else if len(users) != 2 {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("unmatched user amount, want: 2, get: %d", len(users))))
		return
	}

	users, err = kits.SortUsersByUserID(users, []string{operatorID, userID})
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	if users[1].IsLocked {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError("user is already locked"))
		return
	} else if users[0].Permission <= users[1].Permission ||
		users[0].Permission < system_config.GetConfiguration().ARankAdminPermission {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("permission denied, operator: %d, user: %d, need: %d",
			users[0].Permission, users[1].Permission, system_config.GetConfiguration().SRankAdminPermission)))
		return
	}

	users[1].IsLocked = true

	err = dao.GetUser().UpdateColumnsByUserID(users[1], model.User_IsLocked)
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	resData := &struct {
		IsSuccess bool `json:"isSuccess"`
	}{
		IsSuccess: true,
	}

	_, _ = fmt.Fprintln(w, shttp.Response(resData))

	return
}

func unlockUser(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	userID := r.PostFormValue("userID")

	if len(operatorID) < 1 || len(userID) < 1 || operatorID == userID {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, user id: %s", operatorID, userID)))
		return
	}

	users, err := dao.GetUser().Query(model.User_UserID+" in (?)", pg.In([]string{operatorID, userID}))
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	} else if len(users) != 2 {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("unmatched user amount, want 2, get: %d", len(users))))
		return
	}

	users, err = kits.SortUsersByUserID(users, []string{operatorID, userID})
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	if !users[1].IsLocked {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError("user is already unlocked"))
		return
	} else if users[0].Permission <= users[1].Permission ||
		users[0].Permission < system_config.GetConfiguration().ARankAdminPermission {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("permission denied, operator: %d, user: %d, need: %d",
			users[0].Permission, users[1].Permission, system_config.GetConfiguration().SRankAdminPermission)))
		return
	}

	users[1].IsLocked = false

	err = dao.GetUser().UpdateColumnsByUserID(users[1], model.User_IsLocked)
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	resData := &struct {
		IsSuccess bool `json:"isSuccess"`
	}{
		IsSuccess: true,
	}

	_, _ = fmt.Fprintln(w, shttp.Response(resData))

	return
}

func modifyUserPermission(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	userID := r.PostFormValue("userID")
	permissionInt, err := strconv.Atoi(r.PostFormValue("permission"))

	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	if len(operatorID) < 1 || len(userID) < 1 || operatorID == userID {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, user id: %s", operatorID, userID)))
		return
	}

	users, err := dao.GetUser().Query(model.User_UserID+" in (?)", pg.In([]string{operatorID, userID}))
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	} else if len(users) != 2 {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("unmatched user amount, want: 2, get: %d", len(users))))
		return
	}

	users, err = kits.SortUsersByUserID(users, []string{operatorID, userID})
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	permission := uint8(permissionInt)

	if users[0].Permission < system_config.GetConfiguration().SRankAdminPermission ||
		users[1].Permission >= system_config.GetConfiguration().SRankAdminPermission ||
		permission >= system_config.GetConfiguration().SRankAdminPermission {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("permission denied, operator: %d, user: %d, user new: %d",
			users[0].Permission, users[1].Permission, permission)))
		return
	}

	users[1].Permission = permission

	err = dao.GetUser().UpdateColumnsByUserID(users[1], model.User_Permission)
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	resData := &struct {
		IsSuccess bool `json:"isSuccess"`
	}{
		IsSuccess: true,
	}

	_, _ = fmt.Fprintln(w, shttp.Response(resData))

	return
}
