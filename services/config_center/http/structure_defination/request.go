package structure

import "net/http"

// authenticate

type LoginReqParams struct {
    UserName string
    Password string
}

func (p *LoginReqParams) Decode(r *http.Request) string {
    p.UserName = r.PostFormValue(params_UserName)
    p.Password = r.PostFormValue(params_Password)

    return ""
}
