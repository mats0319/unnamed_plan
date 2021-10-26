package http

import (
	"github.com/mats9693/unnamed_plan/admin_data/http/handlers"
	"github.com/mats9693/utils/toy_server/http"
)

var Handlers *mhttp.Handlers

func init() {
	Handlers = mhttp.NewHandlers()

	// user
	Handlers.HandleFunc("/api/login", handlers.Login)

	// cloud file
	Handlers.HandleFunc("/api/cloudFile/listByUploader", handlers.ListCloudFileByUploader)
	Handlers.HandleFunc("/api/cloudFile/listPublic", handlers.ListPublicCloudFile)

	// thinking note
	Handlers.HandleFunc("/api/thinkingNote/listByWriter", handlers.ListThinkingNoteByWriter)
	Handlers.HandleFunc("/api/thinkingNote/listPublic", handlers.ListPublicThinkingNote)
}
