package api

const URI_Register = "/register"

type RegisterReq struct {
	UserName string `json:"user_name"` // unique, also default nickname
	Password string `json:"password"`  // hex(sha256('text'))
}

type RegisterRes struct {
}

const URI_Login = "/login"

type LoginReq struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type LoginRes struct {
	UserName   string `json:"user_name"`
	Nickname   string `json:"nickname"`
	IsAdmin    bool   `json:"is_admin"`
	HasTOTPKey bool   `json:"has_totp_key"`
	EnableMFA  bool   `json:"enable_mfa"` // 为true时上方字段为空，为false时下方字段为空
	MFAToken   string `json:"mfa_token"`
}

const URI_LoginMFA = "/login-mfa"

type LoginMFAReq struct {
	MFAToken string `json:"mfa_token"`
	TOTPCode string `json:"totp_code"`
}

type LoginMFARes struct {
	UserName string `json:"user_name"`
	Nickname string `json:"nickname"`
	IsAdmin  bool   `json:"is_admin"`
}

const URI_ListUser = "/user/list"

type User struct {
	UserName   string `json:"user_name"` // login name, can't modify
	Nickname   string `json:"nickname"`  // display name
	CreatedAt  int64  `json:"created_at"`
	IsAdmin    bool   `json:"is_admin"`
	EnableMFA  bool   `json:"enable_mfa"`
	HasTOTPKey bool   `json:"has_totp_key"`
	LastLogin  int64  `json:"last_login"` // timestamp, unit: milli
}

type ListUserReq struct {
	Page Pagination `json:"page"`
}

type ListUserRes struct {
	Count int64   `json:"count"` // 符合查询条件的用户总数
	Users []*User `json:"users"`
}

const URI_ModifyUser = "/user/modify"

// ModifyUserReq string类型的属性为空，视为不修改对应字段
type ModifyUserReq struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type ModifyUserRes struct {
}

const URI_NewTOTPKey = "/totp-key/new"

type NewTOTPKeyReq struct {
}

type NewTOTPKeyRes struct {
	TOTPKey string `json:"totp_key"`
}

const URI_SetMFAStatus = "/mfa/set-status"

type SetMFAStatusReq struct {
	EnableMFA       bool   `json:"enable_mfa"`
	ApplyNewKeyFlag bool   `json:"apply_new_key_flag"` // 是否申请了新的totp key
	TOTPCode        string `json:"totp_code"`
}

type SetMFAStatusRes struct {
}
