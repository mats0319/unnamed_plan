package model

const (
	Common_ID          = "id"
	Common_CreatedTime = "created_time"
	Common_UpdateTime  = "update_time"
)

const (
	User_UserID     = "user_id"
	User_UserName   = "user_name"
	User_Nickname   = "nickname"
	User_Password   = "password"
	User_Salt       = "salt"
	User_IsLocked   = "is_locked"
	User_Permission = "permission"
	User_CreatedBy  = "created_by"
)

const (
	CloudFile_FileID           = "file_id"
	CloudFile_UploadedBy       = "uploaded_by"
	CloudFile_FileName         = "file_name"
	CloudFile_ExtensionName    = "extension_name"
	CloudFile_LastModifiedTime = "last_modified_time"
	CloudFile_FileSize         = "file_size"
	CloudFile_IsPublic         = "is_public"
	CloudFile_IsDeleted        = "is_deleted"
)

const (
	ThinkingNote_NoteID    = "note_id"
	ThinkingNote_WriteBy   = "write_by"
	ThinkingNote_Topic     = "topic"
	ThinkingNote_Content   = "content"
	ThinkingNote_IsPublic  = "is_public"
	ThinkingNote_IsDeleted = "is_deleted"
)
