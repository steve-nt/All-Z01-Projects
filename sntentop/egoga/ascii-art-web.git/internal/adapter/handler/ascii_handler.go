package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"platform.zone01.gr/git/santonop/SampleAsciiWeb/internal/domain"
	"platform.zone01.gr/git/santonop/SampleAsciiWeb/internal/usecase"
)

type AsciiHandler struct {
	usecase *usecase.AsciiUsecase
}

func NewAsciiHandler(usecase *usecase.AsciiUsecase) *AsciiHandler {
	return &AsciiHandler{usecase: usecase}
}

func (h *AsciiHandler) GenerateAsciiAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var request domain.ASCIITextRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := h.usecase.ConvertTextToAscii(&request)
	if err != nil {
		log.Println("Error generating ASCII art:", err)
		response := domain.AsciiTextResponse{Error: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := domain.AsciiTextResponse{AsciiArt: result}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AsciiHandler) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read the ASCII text from the request body
	err := r.ParseForm()
	if err != nil {
		log.Printf("Failed to parse form: %v", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	asciiArt := r.FormValue("ascii")
	if asciiArt == "" {
		http.Error(w, "No ASCII art provided", http.StatusBadRequest)
		return
	}

	// Set response headers for file download
	w.Header().Set("Content-Disposition", "attachment; filename=ascii_art.txt")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(asciiArt)); err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Failed to write ASCII art to response", http.StatusInternalServerError)
	}
}

func (h *AsciiHandler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 Not Found", http.StatusNotFound)
}
