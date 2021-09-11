package http

import (
	"fmt"
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

	type UserRes struct {
		UserID     string `json:"userID"`
		Nickname   string `json:"nickname"`
		Permission uint8  `json:"permission"`
	}

	var usersRes []*UserRes
	for i := range users {
		usersRes = append(usersRes, &UserRes{
			UserID:     users[i].UserID,
			Nickname:   users[i].Nickname,
			Permission: users[i].Permission,
		})
	}

	resData := &struct {
		Total int        `json:"total"`
		Users []*UserRes `json:"users"`
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
