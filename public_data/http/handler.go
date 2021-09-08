package http

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	h2 "github.com/mats9693/unnamed_plan/admin_data/http"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // todo: for test, write in init, handle together

	name := r.PostFormValue("userName")
	pwd := r.PostFormValue("password")

	users, err := dao.User.Query("name = ?", name)
	if err != nil {
		_, _ = fmt.Fprintln(w, h2.ResponseWithError(err.Error()))
		return
	} else if len(users) != 1 {
		_, _ = fmt.Fprintln(w, h2.ResponseWithError(fmt.Sprintf("invalid data, want %d, get %d.", 1, len(users))))
		return
	}

	user := users[0]
	if user.Password != pwd {
		_, _ = fmt.Fprintln(w, h2.ResponseWithError("invalid account or password"))
		return
	}

	resData := &struct {
		UserID     string `json:"userID"`
		Permission uint8  `json:"permission"`
	}{
		UserID:     user.UserID,
		Permission: user.Permission,
	}

	_, _ = fmt.Fprintln(w, h2.Response(resData))

	return
}
