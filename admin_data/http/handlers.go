package http

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/unnamed_plan/shared/go/config"
	"github.com/mats9693/unnamed_plan/shared/go/http"
	"net/http"
	"strconv"
)

func login(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	name := r.PostFormValue("userName")
	password := r.PostFormValue("password")

	user, err := dao.GetUser().QueryOne("name = ?", name)
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	if user.Password != password {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError("invalid account or password"))
		return
	}

	resData := &struct {
		UserID     string `json:"userID"`
		UserName   string `json:"userName"`
		Permission uint8  `json:"permission"`
	}{
		UserID:     user.UserID,
		UserName:   user.Name,
		Permission: user.Permission,
	}

	_, _ = fmt.Fprintln(w, shttp.Response(resData))

	return
}

func listUser(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	operator, err := dao.GetUser().QueryOne("user_id = ?", operatorID)
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	pageSize, err := strconv.Atoi(r.PostFormValue("pageSize"))
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}
	pageNum, err := strconv.Atoi(r.PostFormValue("pageNum"))
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	users, count, err := dao.GetUser().QueryPage(pageSize, pageNum, "permission <= ?", operator.Permission)
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
		if users[i].UserID != operatorID {
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

func createUser(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	permissionInt, err := strconv.Atoi(r.PostFormValue("permission"))
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	permission := uint8(permissionInt)

	operator, err := dao.GetUser().QueryOne("user_id = ?", operatorID)
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	} else if operator.Permission < config.GetConfiguration().CreateUserPermission ||
		operator.Permission <= permission {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("permission denied, you %d, want create %d.",
			operator.Permission, permission)))
		return
	}

	name := r.PostFormValue("userName")
	password := r.PostFormValue("password")

	err = dao.GetUser().Insert(&model.User{
		Name:       name,
		Nickname:   name,
		Password:   password,
		Permission: permission,
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
	if operatorID == userID {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError("invalid params, operator id is equal to user id"))
		return
	}

	users, err := dao.GetUser().Query("user_id in (?)", pg.In([]string{operatorID, userID}))
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	} else if len(users) != 2 {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("unmatched user amount, want 2, get %d", len(users))))
		return
	}

	users, err = sortUsersByUserID(users, []string{operatorID, userID})
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	if users[1].IsLocked {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError("user is already locked"))
		return
	} else if users[0].Permission <= users[1].Permission ||
		users[0].Permission < config.GetConfiguration().CreateUserPermission { // todo: use lock permission
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("permission denied, you %d, user %d, lock need %d",
			users[0].Permission, users[1].Permission, config.GetConfiguration().CreateUserPermission))) // todo
		return
	}

	users[1].IsLocked = true
	if err = dao.GetUser().UpdateColumnsByUserID(users[1], "is_locked"); err != nil {
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
	if operatorID == userID {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError("invalid params, operator id is equal to user id"))
		return
	}

	users, err := dao.GetUser().Query("user_id in (?)", pg.In([]string{operatorID, userID}))
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	} else if len(users) != 2 {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("unmatched user amount, want 2, get %d", len(users))))
		return
	}

	users, err = sortUsersByUserID(users, []string{operatorID, userID})
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	if !users[1].IsLocked {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError("user is already unlocked"))
		return
	} else if users[0].Permission <= users[1].Permission ||
		users[0].Permission < config.GetConfiguration().CreateUserPermission { // todo: use lock permission
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("permission denied, you %d, user %d, lock need %d",
			users[0].Permission, users[1].Permission, config.GetConfiguration().CreateUserPermission))) // todo
		return
	}

	users[1].IsLocked = false
	if err = dao.GetUser().UpdateColumnsByUserID(users[1], "is_locked"); err != nil {
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
	if operatorID == userID {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError("invalid params, operator id is equal to user id"))
		return
	}

	users, err := dao.GetUser().Query("user_id in (?)", pg.In([]string{operatorID, userID}))
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	} else if len(users) != 2 {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("unmatched user amount, want 2, get %d", len(users))))
		return
	}

	users, err = sortUsersByUserID(users, []string{operatorID, userID})
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	permissionInt, err := strconv.Atoi(r.PostFormValue("permission"))
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	permission := uint8(permissionInt)

	if users[0].Permission < config.GetConfiguration().CreateUserPermission ||
		users[1].Permission >= config.GetConfiguration().CreateUserPermission ||
		permission >= config.GetConfiguration().CreateUserPermission { // todo: use lock permission
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("permission denied, you %d, user %d, lock need %d",
			users[0].Permission, users[1].Permission, config.GetConfiguration().CreateUserPermission))) // todo
		return
	}

	users[1].Permission = permission
	if err = dao.GetUser().UpdateColumnsByUserID(users[1], "permission"); err != nil {
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

func modifyUserInfo(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	userID := r.PostFormValue("userID")
	if operatorID != userID {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError("invalid params, operator id is not equal to user id"))
		return
	}

	user, err := dao.GetUser().QueryOne("user_id = ?", userID)
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	user.Nickname = r.PostFormValue("nickname")
	user.Password = r.PostFormValue("password")

	if err = dao.GetUser().UpdateColumnsByUserID(user, "nickname", "password"); err != nil {
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
