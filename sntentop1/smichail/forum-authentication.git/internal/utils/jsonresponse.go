package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// JsonResponse sends a JSON response with the given data and status code
// This utility function standardizes JSON responses across all handlers
func JsonResponse(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", err)
		// Fallback to a simple error response
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(jsonData); err != nil {
		log.Printf("Failed to write JSON response: %v", err)
	}
}

// ErrorResponse sends a standardized JSON error response
// This ensures consistent error format across the API
func ErrorResponse(w http.ResponseWriter, message string, status int) {
	response := map[string]any{
		"success": false,
		"message": message,
	}
	JsonResponse(w, response, status)
}

// SuccessResponse sends a standardized JSON success response
// func SuccessResponse(w http.ResponseWriter, data any, statusCode int) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
//
// 	response := map[string]any{
// 		"success": true,
// 		"status":  statusCode,
// 	}
//
// 	// If data is a string, treat it as a message
// 	if message, ok := data.(string); ok {
// 		response["message"] = message
// 	} else {
// 		response["data"] = data
// 	}
//
// 	json.NewEncoder(w).Encode(response)
// }
