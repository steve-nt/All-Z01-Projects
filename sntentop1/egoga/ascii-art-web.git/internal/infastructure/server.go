package infrastructure

import (
	"log"
	"net/http"

	"platform.zone01.gr/git/santonop/SampleAsciiWeb/internal/adapter/handler"
)

func StartServer(handler *handler.AsciiHandler) {
	router := SetupRouter(handler)
	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
