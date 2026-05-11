package infrastructure

import (
	"net/http"

	"platform.zone01.gr/git/santonop/SampleAsciiWeb/internal/adapter/handler"
)

func SetupRouter(handler *handler.AsciiHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ascii-art", handler.GenerateAsciiAPI)
	mux.HandleFunc("/", handler.NotFoundHandler) // Catch-all for undefined routes
	mux.HandleFunc("/download", handler.DownloadHandler)
	return mux
}
