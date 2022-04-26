package mconst

const (
	LogFileName    = "log.log"
	ConfigFileName = "config.json"

	EmptyHTTPRes = ""

	// uid
	UID_Config             = "612fbb57-c44c-4b54-8188-13d1dd598306"
	UID_DB_Dev             = "658e06f7-71d5-4ada-b715-8c1a4489e5d2"
	UID_DB_Production      = "658e06f7-71d5-4ada-b715-8c1a4489e5d3"
	UID_Gateway_HTTP       = "3b839c1f-9f1e-474b-b3da-7b5e9bc792ec"
	UID_Gateway_RPC_Client = "1cd10cb8-ecf5-4855-a886-76b148ed104a"

	UID_Gateway_Service    = "84d1fecc-3be9-439e-8144-209ffc00a975"
	UID_User_Service       = "eafbda7d-c951-4fc9-8b45-8c90189c1e36"
	UID_Cloud_File_Service = "1b5ab1d2-de6d-4377-9a4e-a184b24d1a0f"
	UID_Note_Service       = "23d062e4-3c36-45f0-9e1c-3f339742903b"
	UID_Task_Service       = "a4802e2b-113b-4132-b125-ca5f97239a6e"

	// db
	DB_PostgreSQL = "postgresql"

	Error_UnsupportedDB = "unsupported db: "

	// http
	HTTP_SourceSign = "Unnamed-Plan-Source"

	HTTP_MultiLoginUserIDSign = "Unnamed-Plan-User"
	HTTP_MultiLoginTokenSign  = "Unnamed-Plan-Token"

	Error_UnknownSource       = "unknown request source: "
	Error_UnknownURI          = "unknown request URI: "
	Error_LoadLoginInfoFailed = "invalid login info"
	Error_InvalidToken        = "invalid token"
	Error_InvalidTokenTimeout = "invalid token: timeout"
	Error_InvalidParams       = "invalid param(s)"
)

const (
	Error_InvalidAccountOrPassword = "invalid account or password"
	Error_PermissionDenied         = "permission denied"
	Error_NoValidModification      = "not any valid modification"

	Error_UserAlreadyLocked   = "user already locked"
	Error_UserAlreadyUnlocked = "user already unlocked"
	Error_FileAlreadyDeleted  = "file already deleted"
	Error_NoteAlreadyDeleted  = "note already deleted"
	Error_ModifyOthersNote    = "not allowed to modify others' note"
	Error_ModifyOthersTask    = "not allowed to modify others' task"
)
