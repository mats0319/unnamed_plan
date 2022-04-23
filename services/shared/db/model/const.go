package model

const (
	Common_ID             = "id"
	Common_CreatedTime    = "created_time"
	Common_UpdateTime     = "update_time"
	Common_OptimisticLock = "optimistic_lock"
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
	Note_WriteBy   = "write_by"
	Note_Topic     = "topic"
	Note_Content   = "content"
	Note_IsPublic  = "is_public"
	Note_IsDeleted = "is_deleted"
)

const (
	Task_TaskName    = "task_name"
	Task_PostedBy    = "posted_by"
	Task_Description = "description"
	Task_PreTaskIDs  = "pre_task_ids"
	Task_Status      = "status"
)

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
	ConfigItem_UsedIn         = "used_in"
)
