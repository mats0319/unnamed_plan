package http

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/unnamed_plan/admin_data/kits"
	"github.com/mats9693/utils/toy_server/http"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	name := r.PostFormValue("userName")
	password := r.PostFormValue("password")

	user, err := dao.GetUser().QueryOne(model.User_UserName+" = ?", name)
	if err != nil {
		_, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
		return
	}

	if kits.VerifyUserPassword(user.Password, password, user.Salt) {
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
