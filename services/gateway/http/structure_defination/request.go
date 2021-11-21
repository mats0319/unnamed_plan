package structure

import (
	"github.com/mats9693/unnamed_plan/services/utils"
	"mime/multipart"
	"net/http"
	"strconv"
)

// user

type LoginReqParams struct {
	UserName string
	Password string
}

func (p *LoginReqParams) Decode(r *http.Request) string {
	p.UserName = r.PostFormValue(params_UserName)
	p.Password = r.PostFormValue(params_Password)

	return ""
}

type ListUserReqParams struct {
	OperatorID string
	PageSize   int
	PageNum    int
}

func (p *ListUserReqParams) Decode(r *http.Request) string {
	p.OperatorID = r.PostFormValue(params_OperatorID)
	pageSize, err := strconv.Atoi(r.PostFormValue(params_PageSize))
	pageNum, err2 := strconv.Atoi(r.PostFormValue(params_PageNum))

	if err != nil || err2 != nil {
		return utils.ErrorsToString(err, err2)
	}

	p.PageSize = pageSize
	p.PageNum = pageNum

	return ""
}

type CreateUserReqParams struct {
	OperatorID string
	UserName   string
	Password   string
	Permission uint8
}

func (p *CreateUserReqParams) Decode(r *http.Request) string {
	p.OperatorID = r.PostFormValue(params_OperatorID)
	p.UserName = r.PostFormValue(params_UserName)
	p.Password = r.PostFormValue(params_Password)
	permissionInt, err := strconv.Atoi(r.PostFormValue(params_Permission))

	if err != nil {
		return err.Error()
	}

	p.Permission = uint8(permissionInt)

	return ""
}

type LockUserReqParams struct {
	OperatorID string
	UserID     string
}

func (p *LockUserReqParams) Decode(r *http.Request) string {
	p.OperatorID = r.PostFormValue(params_OperatorID)
	p.UserID = r.PostFormValue(params_UserID)

	return ""
}

type UnlockUserReqParams struct {
	OperatorID string
	UserID     string
}

func (p *UnlockUserReqParams) Decode(r *http.Request) string {
	p.OperatorID = r.PostFormValue(params_OperatorID)
	p.UserID = r.PostFormValue(params_UserID)

	return ""
}

type ModifyUserInfoReqParams struct {
	OperatorID string
	UserID     string
	CurrPwd    string
	Nickname   string
	Password   string
}

func (p *ModifyUserInfoReqParams) Decode(r *http.Request) string {
	p.OperatorID = r.PostFormValue(params_OperatorID)
	p.UserID = r.PostFormValue(params_UserID)
	p.CurrPwd = r.PostFormValue(params_CurrPwd)
	p.Nickname = r.PostFormValue(params_Nickname)
	p.Password = r.PostFormValue(params_Password)

	return ""
}

type ModifyUserPermissionReqParams struct {
	OperatorID string
	UserID     string
	Permission uint8
}

func (p *ModifyUserPermissionReqParams) Decode(r *http.Request) string {
	p.OperatorID = r.PostFormValue(params_OperatorID)
	p.UserID = r.PostFormValue(params_UserID)
	permissionInt, err := strconv.Atoi(r.PostFormValue(params_Permission))

	if err != nil {
		return err.Error()
	}

	p.Permission = uint8(permissionInt)

	return ""
}

// cloud file

type ListCloudFileByUploaderReqParams struct {
	OperatorID string
	PageSize   int
	PageNum    int
}

func (p *ListCloudFileByUploaderReqParams) Decode(r *http.Request) string {
	p.OperatorID = r.PostFormValue(params_OperatorID)
	pageSize, err := strconv.Atoi(r.PostFormValue(params_PageSize))
	pageNum, err2 := strconv.Atoi(r.PostFormValue(params_PageNum))

	if err != nil || err2 != nil {
		return utils.ErrorsToString(err, err2)
	}

	p.PageSize = pageSize
	p.PageNum = pageNum

	return ""
}

type ListPublicCloudFileReqParams struct {
	OperatorID string
	PageSize   int
	PageNum    int
}

func (p *ListPublicCloudFileReqParams) Decode(r *http.Request) string {
	p.OperatorID = r.PostFormValue(params_OperatorID)
	pageSize, err := strconv.Atoi(r.PostFormValue(params_PageSize))
	pageNum, err2 := strconv.Atoi(r.PostFormValue(params_PageNum))

	if err != nil || err2 != nil {
		return utils.ErrorsToString(err, err2)
	}

	p.PageSize = pageSize
	p.PageNum = pageNum

	return ""
}

type UploadCloudFileReqParams struct {
	OperatorID       string
	FileName         string
	ExtensionName    string
	LastModifiedTime int
	IsPublic         bool
	File             []byte
	FileHeader       *multipart.FileHeader
}

func (p *UploadCloudFileReqParams) Decode(r *http.Request) string {
	p.OperatorID = r.PostFormValue(params_OperatorID)
	p.FileName = r.PostFormValue(params_FileName)
	p.ExtensionName = r.PostFormValue(params_ExtensionName)
	lastModifiedTime, err := strconv.Atoi(r.PostFormValue(params_LastModifiedTime))
	isPublic, err2 := utils.StringToBool(r.PostFormValue(params_IsPublic))
	file, fileHeader, err3 := r.FormFile(params_File)

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

	p.LastModifiedTime = lastModifiedTime
	p.IsPublic = isPublic
	p.File = fileContent
	p.FileHeader = fileHeader

	return ""
}

type ModifyCloudFileReqParams struct {
	// required
	OperatorID string
	FileID     string
	Password   string

	// option
	FileName      string
	ExtensionName string
	IsPublic      bool
	File          []byte

	// required if 'file' is not null
	FileHeader       *multipart.FileHeader
	LastModifiedTime int
}

func (p *ModifyCloudFileReqParams) Decode(r *http.Request) string {
	p.OperatorID = r.PostFormValue(params_OperatorID)
	p.FileID = r.PostFormValue(params_FileID)
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
		p.FileHeader = fileHeader
		p.LastModifiedTime = lastModifiedTime
	}

	p.IsPublic = isPublic

	return ""
}

type DeleteCloudFileReqParams struct {
	OperatorID string
	Password   string
	FileID     string
}

func (p *DeleteCloudFileReqParams) Decode(r *http.Request) string {
	p.OperatorID = r.PostFormValue(params_OperatorID)
	p.Password = r.PostFormValue(params_Password)
	p.FileID = r.PostFormValue(params_FileID)

	return ""
}

// thinking note

type ListThinkingNoteByWriterReqParams struct {
	OperatorID string
	PageSize   int
	PageNum    int
}

func (p *ListThinkingNoteByWriterReqParams) Decode(r *http.Request) string {
	p.OperatorID = r.PostFormValue(params_OperatorID)
	pageSize, err := strconv.Atoi(r.PostFormValue(params_PageSize))
	pageNum, err2 := strconv.Atoi(r.PostFormValue(params_PageNum))

	if err != nil || err2 != nil {
		return utils.ErrorsToString(err, err2)
	}

	p.PageSize = pageSize
	p.PageNum = pageNum

	return ""
}

type ListPublicThinkingNoteReqParams struct {
	OperatorID string
	PageSize   int
	PageNum    int
}

func (p *ListPublicThinkingNoteReqParams) Decode(r *http.Request) string {
	p.OperatorID = r.PostFormValue(params_OperatorID)
	pageSize, err := strconv.Atoi(r.PostFormValue(params_PageSize))
	pageNum, err2 := strconv.Atoi(r.PostFormValue(params_PageNum))

	if err != nil || err2 != nil {
		return utils.ErrorsToString(err, err2)
	}

	p.PageSize = pageSize
	p.PageNum = pageNum

	return ""
}

type CreateThinkingNoteReqParams struct {
	OperatorID string
	Topic      string
	Content    string
	IsPublic   bool
}

func (p *CreateThinkingNoteReqParams) Decode(r *http.Request) string {
	p.OperatorID = r.PostFormValue(params_OperatorID)
	p.Topic = r.PostFormValue(params_Topic)
	p.Content = r.PostFormValue(params_Content)
	isPublic, err := utils.StringToBool(r.PostFormValue(params_IsPublic))

	if err != nil {
		return err.Error()
	}

	p.IsPublic = isPublic

	return ""
}

type ModifyThinkingNoteReqParams struct {
	OperatorID string
	NoteID     string
	Password   string
	Topic      string
	Content    string
	IsPublic   bool
}

func (p *ModifyThinkingNoteReqParams) Decode(r *http.Request) string {
	p.OperatorID = r.PostFormValue("operatorID")
	p.NoteID = r.PostFormValue("noteID")
	p.Password = r.PostFormValue("password")
	p.Topic = r.PostFormValue("topic")
	p.Content = r.PostFormValue("content")
	isPublic, err := utils.StringToBool(r.PostFormValue("isPublic"))

	if err != nil {
		return err.Error()
	}

	p.IsPublic = isPublic

	return ""
}

type DeleteThinkingNoteReqParams struct {
	OperatorID string
	Password   string
	NoteID     string
}

func (p *DeleteThinkingNoteReqParams) Decode(r *http.Request) string {
	p.OperatorID = r.PostFormValue("operatorID")
	p.Password = r.PostFormValue("password")
	p.NoteID = r.PostFormValue("noteID")

	return ""
}
