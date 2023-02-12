package structure

import (
	"github.com/mats9693/unnamed_plan/services/shared/proto/go"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"net/http"
	"strconv"
)

// request params

const (
	params_UserName         = "user_name"
	params_Password         = "password"
	params_OperatorID       = "operator_id"
	params_PageSize         = "page_size"
	params_PageNum          = "page_num"
	params_Permission       = "permission"
	params_UserID           = "user_id"
	params_CurrPwd          = "curr_pwd"
	params_Nickname         = "nickname"
	params_OperatorPassword = "operator_password"

	params_Rule             = "rule"
	params_FileName         = "file_name"
	params_ExtensionName    = "extension_name"
	params_LastModifiedTime = "lastModified_time"
	params_IsPublic         = "is_public"
	params_File             = "file"
	params_FileID           = "file_id"
)

// user

type LoginReqParams rpc_impl.User_LoginReq

func (p *LoginReqParams) Decode(r *http.Request) {
	p.UserName = r.PostFormValue(params_UserName)
	p.Password = r.PostFormValue(params_Password)
}

type ListUserReqParams rpc_impl.User_ListReq

func (p *ListUserReqParams) Decode(r *http.Request) string {
	p.OperatorId = r.PostFormValue(params_OperatorID)
	pageSize, err := strconv.Atoi(r.PostFormValue(params_PageSize))
	pageNum, err2 := strconv.Atoi(r.PostFormValue(params_PageNum))

	if err != nil || err2 != nil {
		return utils.ErrorsToString(err, err2)
	}

	p.Page = &rpc_impl.Pagination{
		PageSize: uint32(pageSize),
		PageNum:  uint32(pageNum),
	}

	return ""
}

type CreateUserReqParams rpc_impl.User_CreateReq

func (p *CreateUserReqParams) Decode(r *http.Request) string {
	p.OperatorId = r.PostFormValue(params_OperatorID)
	p.UserName = r.PostFormValue(params_UserName)
	p.Password = r.PostFormValue(params_Password)
	permissionInt, err := strconv.Atoi(r.PostFormValue(params_Permission))
	p.OperatorPassword = r.PostFormValue(params_OperatorPassword)

	if err != nil {
		return err.Error()
	}

	p.Permission = uint32(permissionInt)

	return ""
}

type LockUserReqParams rpc_impl.User_LockReq

func (p *LockUserReqParams) Decode(r *http.Request) {
	p.OperatorId = r.PostFormValue(params_OperatorID)
	p.UserId = r.PostFormValue(params_UserID)
	p.Password = r.PostFormValue(params_Password)
}

type UnlockUserReqParams rpc_impl.User_UnlockReq

func (p *UnlockUserReqParams) Decode(r *http.Request) {
	p.OperatorId = r.PostFormValue(params_OperatorID)
	p.UserId = r.PostFormValue(params_UserID)
	p.Password = r.PostFormValue(params_Password)
}

type ModifyUserInfoReqParams rpc_impl.User_ModifyInfoReq

func (p *ModifyUserInfoReqParams) Decode(r *http.Request) {
	p.OperatorId = r.PostFormValue(params_OperatorID)
	p.UserId = r.PostFormValue(params_UserID)
	p.CurrPwd = r.PostFormValue(params_CurrPwd)
	p.Nickname = r.PostFormValue(params_Nickname)
	p.Password = r.PostFormValue(params_Password)
}

type ModifyUserPermissionReqParams rpc_impl.User_ModifyPermissionReq

func (p *ModifyUserPermissionReqParams) Decode(r *http.Request) string {
	p.OperatorId = r.PostFormValue(params_OperatorID)
	p.UserId = r.PostFormValue(params_UserID)
	permissionInt, err := strconv.Atoi(r.PostFormValue(params_Permission))
	p.Password = r.PostFormValue(params_Password)

	if err != nil {
		return err.Error()
	}

	p.Permission = uint32(permissionInt)

	return ""
}

// cloud file

type ListCloudFileReqParams rpc_impl.CloudFile_ListReq

func (p *ListCloudFileReqParams) Decode(r *http.Request) string {
	rule, err := strconv.Atoi(r.PostFormValue(params_Rule))
	p.OperatorId = r.PostFormValue(params_OperatorID)
	pageSize, err2 := strconv.Atoi(r.PostFormValue(params_PageSize))
	pageNum, err3 := strconv.Atoi(r.PostFormValue(params_PageNum))

	if err != nil || err2 != nil || err3 != nil {
		return utils.ErrorsToString(err, err2, err3)
	}

	p.Rule = rpc_impl.CloudFile_ListRule(rule)
	p.Page = &rpc_impl.Pagination{
		PageSize: uint32(pageSize),
		PageNum:  uint32(pageNum),
	}

	return ""
}

type UploadCloudFileReqParams rpc_impl.CloudFile_UploadReq

func (p *UploadCloudFileReqParams) Decode(r *http.Request) string {
	p.OperatorId = r.PostFormValue(params_OperatorID)
	p.FileName = r.PostFormValue(params_FileName)
	p.ExtensionName = r.PostFormValue(params_ExtensionName)
	lastModifiedTime, err := strconv.Atoi(r.PostFormValue(params_LastModifiedTime))
	isPublic, err2 := utils.StringToBool(r.PostFormValue(params_IsPublic))
	file, fileHeader, err3 := r.FormFile(params_File)
	p.Password = r.PostFormValue(params_Password)

	if err != nil || err2 != nil || err3 != nil {
		return utils.ErrorsToString(err, err2, err3)
	}

	fileContent := make([]byte, fileHeader.Size) // require enough length before read
	_, err = file.Read(fileContent)
	if err != nil {
		return err.Error()
	}
	defer func() {
		_ = file.Close()
	}()

	p.LastModifiedTime = int64(lastModifiedTime)
	p.IsPublic = isPublic
	p.File = fileContent
	p.FileSize = fileHeader.Size

	return ""
}

type ModifyCloudFileReqParams rpc_impl.CloudFile_ModifyReq

func (p *ModifyCloudFileReqParams) Decode(r *http.Request) string {
	p.OperatorId = r.PostFormValue(params_OperatorID)
	p.FileId = r.PostFormValue(params_FileID)
	p.Password = r.PostFormValue(params_Password)
	p.FileName = r.PostFormValue(params_FileName)
	p.ExtensionName = r.PostFormValue(params_ExtensionName)
	isPublic, err := utils.StringToBool(r.PostFormValue(params_IsPublic))
	file, fileHeader, err2 := r.FormFile(params_File)
	lastModifiedTime, err3 := strconv.Atoi(r.PostFormValue(params_LastModifiedTime))

	if err != nil || (err2 != nil && err2 != http.ErrMissingFile) || err3 != nil {
		return utils.ErrorsToString(err, err2, err3)
	}

	if err2 == nil {
		fileContent := make([]byte, fileHeader.Size) // require enough length before read
		_, err = file.Read(fileContent)
		if err != nil {
			return err.Error()
		}
		defer func() {
			_ = file.Close()
		}()

		p.File = fileContent
		p.FileSize = fileHeader.Size
		p.LastModifiedTime = int64(lastModifiedTime)
	}

	p.IsPublic = isPublic

	return ""
}

type DeleteCloudFileReqParams rpc_impl.CloudFile_DeleteReq

func (p *DeleteCloudFileReqParams) Decode(r *http.Request) {
	p.OperatorId = r.PostFormValue(params_OperatorID)
	p.FileId = r.PostFormValue(params_FileID)
	p.Password = r.PostFormValue(params_Password)
}
