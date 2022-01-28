package mconst

const (
	LogFileName    = "log.log"
	ConfigFileName = "config.json"
	ConfigDevLevel = "dev"
	EmptyHTTPRes   = ""

	// uid
	UID_DB   = "658e06f7-71d5-4ada-b715-8c1a4489e5d2"
	UID_HTTP = "3b839c1f-9f1e-474b-b3da-7b5e9bc792ec"

	// db
	DB_PostgreSQL = "postgresql"

	Error_UnsupportedDB = "unsupported db: "

	// http
	HTTP_MultiLoginSourceSign = "Unnamed-Plan-Source"
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

	Error_UserAlreadyLocked        = "user already locked"
	Error_UserAlreadyUnlocked      = "user already unlocked"
	Error_FileAlreadyDeleted       = "file already deleted"
	Error_NoteAlreadyDeleted       = "note already deleted"
	Error_ModifyOthersThinkingNote = "not allowed to modify others' thinking note"
	Error_ModifyOthersTask         = "not allowed to modify others' task"
)
