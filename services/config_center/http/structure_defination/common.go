package structure

// request params
const (
    params_UserName   = "userName"
    params_Password   = "password"
    params_OperatorID = "operatorID"
    params_PageSize   = "pageSize"
    params_PageNum    = "pageNum"
    params_Permission = "permission"
    params_UserID     = "userID"
    params_CurrPwd    = "currPwd"
    params_Nickname   = "nickname"

    params_FileName         = "fileName"
    params_ExtensionName    = "extensionName"
    params_LastModifiedTime = "lastModifiedTime"
    params_IsPublic         = "isPublic"
    params_File             = "file"
    params_FileID           = "fileID"

    params_Topic   = "topic"
    params_Content = "content"
)

// common

type Total struct {
    Total uint32 `json:"total"`
}

// authenticate

type UserID struct {
    UserID string `json:"userID"`
}

type UserName struct {
    UserName string `json:"userName"`
}
