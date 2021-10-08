package http

import (
	"github.com/mats9693/unnamed_plan/admin_data/http/handlers"
	"net/http"
)

func init() {
	// user
	http.HandleFunc("/api/login", handlers.Login)

	// cloud file
	http.HandleFunc("/api/cloudFile/listByUploader", handlers.ListCloudFileByUploader)
	http.HandleFunc("/api/cloudFile/listPublic", handlers.ListPublicCloudFile)

	// thinking note
	http.HandleFunc("/api/thinkingNote/listByWriter", handlers.ListThinkingNoteByWriter)
	http.HandleFunc("/api/thinkingNote/listPublic", handlers.ListPublicThinkingNote)
}
