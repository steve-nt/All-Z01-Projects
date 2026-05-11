package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type RequestData struct {
	Text string `json:"text"`
}

type ResponseData struct {
	Response string `json:"response"`
	Image    string `json:"image"` // Base64-encoded image
}

func main() {
	// Serve the index.html file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("templates", "index.html"))
	})

	// Handle POST requests to /process
	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Parse the JSON request body
		var reqData RequestData
		err := json.NewDecoder(r.Body).Decode(&reqData)
		if err != nil {
			http.Error(w, "Failed to parse request", http.StatusBadRequest)
			return
		}

		// Read and encode the PNG file
		imagePath := filepath.Join("templates", "premium-french-fries-photos-7-png-500x500.png")
		imageData, err := ioutil.ReadFile(imagePath)
		if err != nil {
			http.Error(w, "Failed to read image file", http.StatusInternalServerError)
			return
		}
		encodedImage := base64.StdEncoding.EncodeToString(imageData)

		// Create the response
		respData := ResponseData{
			Response: fmt.Sprintf("You said: %s", reqData.Text),
			Image:    encodedImage,
		}

		// Send JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respData)
	})

	// Start the HTTP server
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
