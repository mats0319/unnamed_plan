package http

import (
    "fmt"
    "github.com/mats9693/unnamed_plan/admin_data/db/dao"
    "github.com/mats9693/unnamed_plan/shared/go/http"
    "net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
    if isDev {
        w.Header().Set("Access-Control-Allow-Origin", "*")
    }

    name := r.PostFormValue("userName")
    password := r.PostFormValue("password")

    users, err := dao.GetUser().Query("name = ?", name) // todo: add unlock required
    if err != nil {
        _, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
        return
    } else if len(users) != 1 {
        _, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("invalid data, want %d, get %d.", 1, len(users))))
        return
    }

    user := users[0]
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
