package http

import (
	"fmt"
	"net/http"

	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/unnamed_plan/admin_data/kits"
	mhttp "github.com/mats9693/utils/toy_server/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	userName := r.PostFormValue("userName")
	password := r.PostFormValue("password")

	user, err := dao.GetUser().QueryOne(model.User_UserName+" = ?", userName)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	if user.Password != kits.CalcSHA256(password, user.Salt) {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError("invalid account or password"))
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
