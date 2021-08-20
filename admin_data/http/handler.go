package http

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/public_data/db/dao"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // for test

	name := r.PostFormValue("userName")
	pwd := r.PostFormValue("password")

	user, err := dao.GetUser(name, pwd)
	if err != nil {
		_, _ = fmt.Fprintln(w, responseWithError(err.Error()))
		return
	}

	resData := &struct {
		Permission uint8 `json:"permission"`
	}{
		Permission: user.Permission,
	}

	_, _ = fmt.Fprintln(w, response(resData))

	return
}
