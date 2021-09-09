package http

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/admin_data/config"
	. "github.com/mats9693/unnamed_plan/admin_data/const"
	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	"net/http"
	"strconv"
)

func login(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	name := r.PostFormValue(UserName)
	password := r.PostFormValue(Password)

	users, err := dao.GetUser().Query("name = ?", name) // todo: add unlock required
	if err != nil {
		_, _ = fmt.Fprintln(w, ResponseWithError(err.Error()))
		return
	} else if len(users) != 1 {
		_, _ = fmt.Fprintln(w, ResponseWithError(fmt.Sprintf("invalid data, want %d, get %d.", 1, len(users))))
		return
	}

	user := users[0]
	if user.Password != password {
		_, _ = fmt.Fprintln(w, ResponseWithError("invalid account or password"))
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

	_, _ = fmt.Fprintln(w, Response(resData))

	return
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	userID := r.PostFormValue(UserID)
	permissionInt, err := strconv.Atoi(r.PostFormValue(Permission))
	if err != nil {
		_, _ = fmt.Fprintln(w, ResponseWithError(err.Error()))
		return
	}

	permission := uint8(permissionInt)

	users, err := dao.GetUser().Query("user_id = ?", userID)
	if err != nil {
		_, _ = fmt.Fprintln(w, ResponseWithError(err.Error()))
		return
	} else if len(users) != 1 {
		_, _ = fmt.Fprintln(w, ResponseWithError(fmt.Sprintf("invalid data, want %d, get %d.", 1, len(users))))
		return
	} else if users[0].Permission < config.GetConfiguration().CreateUserPermission ||
		users[0].Permission <= permission {
		_, _ = fmt.Fprintln(w, ResponseWithError(fmt.Sprintf("permission denied, you %d, want create %d.",
			users[0].Permission, permission)))
		return
	}

	name := r.PostFormValue(UserName)
	password := r.PostFormValue(Password)

	err = dao.GetUser().Insert(name, password, permission)
	if err != nil {
		_, _ = fmt.Fprintln(w, ResponseWithError(err.Error()))
		return
	}

	resData := &struct {
		IsSuccess bool `json:"isSuccess"`
	}{
		IsSuccess: true,
	}

	_, _ = fmt.Fprintln(w, Response(resData))

	return
}
