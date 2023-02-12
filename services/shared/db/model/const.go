package model

// config center

const (
	ServiceConfig_ServiceID     = "service_id"
	ServiceConfig_Level         = "level"
	ServiceConfig_ServiceName   = "service_name"
	ServiceConfig_ConfigItemIDs = "config_item_ids"
	ServiceConfig_IsDelete      = "is_delete"
)

const (
	ConfigItem_ConfigItemID   = "config_item_id"
	ConfigItem_ConfigItemName = "config_item_name"
	ConfigItem_ConfigItemTag  = "config_item_tag"
	ConfigItem_ConfigSubItems = "config_sub_items"
	ConfigItem_BeUsed         = "be_used"
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

// common

const (
	Common_ID             = "id"
	Common_CreatedTime    = "created_time"
	Common_UpdateTime     = "update_time"
	Common_OptimisticLock = "optimistic_lock"
)
