package model

const (
	Common_ID          = "id"
	Common_CreatedTime = "created_time"
	Common_UpdateTime  = "update_time"
)

// config center

const (
	Administrator_UserName = "user_name"
	Administrator_Password = "password"
	Administrator_Salt     = "salt"
)

const (
	Version_VersionNum     = "version_num"
	Version_Description    = "description"
	Version_ServicesIDs    = "service_ids"
	Version_ConfigIDs      = "config_ids"
	Version_Configurations = "configurations"
	Version_IsUsing        = "is_using"
	Version_HasUpdate      = "has_update"
)

const (
	Service_ServiceID = "service_id"
	Service_ServiceName = "service_name"
	Service_ConfigIDs = "config_ids"
	Service_IsShadow = "is_shadow"
)

const (
	Config_ConfigID = "config_id"
	Config_ConfigName = "config_name"
	Config_payload = "payload"
	Config_IsShadow = "is_shadow"
)

// services

const (
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
	ThinkingNote_WriteBy   = "write_by"
	ThinkingNote_Topic     = "topic"
	ThinkingNote_Content   = "content"
	ThinkingNote_IsPublic  = "is_public"
	ThinkingNote_IsDeleted = "is_deleted"
)
