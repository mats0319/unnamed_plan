package structure

import "net/http"

type LoginParams struct {
    UserName string
    Password string
}

func DecodeLoginParams(r *http.Request) (*LoginParams, error) {
    params := &LoginParams{}

    params.UserName = r.PostFormValue(params_UserName)
    params.Password = r.PostFormValue(params_Password)

    return params, nil
}
