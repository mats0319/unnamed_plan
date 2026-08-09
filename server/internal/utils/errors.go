package utils

// 每个变量代表一个函数，一个预设了不同参数的`NewError`函数，使用时需要添加`()`
// 错误码为5位十进制数：
// - 业务逻辑错误（http状态码为200）：第一位表示错误类型，2～3位表示功能模块，4～5位表示错误序号
// - 非业务逻辑错误：第1～3位表示http状态码，4～5位表示错误序号
// 使用invalid表示结构错误，使用wrong表示数值错误
var (
	ErrForTest = newError(0, -1, "test error string")

	// change http code error
	ErrInvalidAccessToken  = newError(401, 40101, "Invalid Access Token")
	ErrServerInternalError = newError(500, 50001, "Sever Internal Error")
	ErrDBError             = newError(500, 50002, "Sever Internal Error")

	// params error (1)
	// general (00) / db (01) / middleware (02)
	ErrDeserializeReqParam = newBusinessError(10001, "Deserialize HTTP Request Params Failed")
	ErrInvalidParams       = newBusinessError(10002, "Invalid Params")
	ErrEncrypt             = newBusinessError(10003, "Encrypt Failed")
	ErrDecrypt             = newBusinessError(10004, "Decrypt Failed")

	ErrUserExist    = newBusinessError(10101, "User Already Exist")
	ErrNoteExist    = newBusinessError(10102, "Note Already Exist")
	ErrUserNotFound = newBusinessError(10103, "User Not Found")
	ErrNoteNotFound = newBusinessError(10104, "Note Not Found")

	ErrInvalidToken     = newBusinessError(10201, "Invalid Token")
	ErrWrongTokenHash   = newBusinessError(10202, "Wrong Token Hash")
	ErrDecodeToken      = newBusinessError(10203, "Decode Token Failed")
	ErrDeserializeToken = newBusinessError(10204, "Deserialize Token Failed")
	ErrInvalidTokenType = newBusinessError(10205, "Invalid Token Type")
	ErrTokenExpired     = newBusinessError(10206, "Token Expired")

	// business error (2)
	// general (00) / user (01) / note (02) / game score (03)
	ErrNoChanges        = newBusinessError(20001, "No Changes")
	ErrPermissionDenied = newBusinessError(20002, "Permission Denied")

	ErrInvalidPassword = newBusinessError(20101, "Invalid Password")
	ErrInvalidPwdSalt  = newBusinessError(20102, "Invalid Password Salt")
	ErrInvalidPwdKey   = newBusinessError(20103, "Invalid Password Key")
	ErrWrongPassword   = newBusinessError(20104, "Wrong Password")
	ErrInvalidTOTPCode = newBusinessError(20105, "Invalid TOTP Code")
	ErrInvalidTOTPKey  = newBusinessError(20106, "Invalid TOTP Key")
	ErrWrongTOTPCode   = newBusinessError(20107, "Wrong TOTP Code")
	ErrSamePassword    = newBusinessError(20108, "New Password Can't be Identical to the Old One")
	ErrTOTPKeyNotFound = newBusinessError(20109, "TOTP Key not Found")
	ErrTryTooManyTimes = newBusinessError(20110, "Try Too Many Times")
	ErrTOTPKeyExpired  = newBusinessError(20111, "TOTP Key Expired")

	ErrInvalidGameName = newBusinessError(20301, "Invalid Game Name")
)

// 函数返回函数而不是实例，可以避免多处使用同一变量会继承历史数据的问题，详见测试代码
func newError(httpCode int, code int, detail string) func() *Error {
	return func() *Error { return NewError(httpCode, code, detail) }
}

func newBusinessError(code int, detail string) func() *Error {
	return newError(200, code, detail)
}
