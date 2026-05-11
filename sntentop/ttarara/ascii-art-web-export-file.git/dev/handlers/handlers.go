package handlers

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func PgHome(w http.ResponseWriter, r *http.Request) {
	data := PrepareTemplateData("HOME", "home-page")
	log.Println("requesting_page: HOME") // na ginei me LogEventsRecord
	RenderTemplateFiles(w, "PgHome.html", data)
	log.Printf("Rendering template %v successfully\n", data["Title"])
}

func PgProject(w http.ResponseWriter, r *http.Request) {
	data := PrepareTemplateData("PROJECT", "project-page")
	log.Println("requesting_page: PROJECT") // na ginei me LogEventsRecord
	RenderTemplateFiles(w, "PgProject.html", data)
	log.Printf("Rendering template %v successfully\n", data["Title"])
}

func PgTeam(w http.ResponseWriter, r *http.Request) {
	data := PrepareTemplateData("TEAM", "team-page")
	log.Println("requesting_page: TEAM") // na ginei me LogEventsRecord
	RenderTemplateFiles(w, "PgTeam.html", data)
	log.Printf("Rendering template %v successfully\n", data["Title"])
}

func PgConverter(w http.ResponseWriter, r *http.Request) {
	data := PrepareTemplateData("CONVERTER", "converter-page")
	if r.Method == http.MethodPost {

		text := r.FormValue("ascii-data") // Handle form submission
		banner := r.FormValue("banner")   // Handle form submission

		final, renderStatus := inputValidation(text) // Validate and process input
		log.Println("Render Status:", renderStatus, "Rendered Text:", final)
		if renderStatus == 400 {
			log.Println("400 Error", "Invalid Input: "+text)
			handleBadRequest(w, "Invalid Input: "+text)
			return
		}

		if renderStatus != 200 {
			data["Error"] = final
		} else {
			output, statusCode := asciiArt(final, banner)
			if statusCode == 500 {
				log.Println("500 Error", "Failed to read font file for font: "+banner)
				handleServerError(w, "Failed to generate ASCII art")
				return
			}

			if statusCode != 200 {
				data["Error"] = "Failed to generate ASCII art. Please check your input and try again."
			} else {
				data["First"] = output // Set the ASCII Art result
				fpath, err := ExportToFile(output)
				if err != nil {
					log.Printf("Error: %v", err)
				} else {
					log.Printf("File successfully exported to: %s", fpath)
				}

			}
		}
	}
	log.Println("requesting-page: CONVERTER") // na ginei me LogEventsRecord
	RenderTemplateFiles(w, "PgConverter.html", data)
	log.Printf("Rendering template %v successfully\n", data["Title"])
	log.Printf("The output:\n %v\n", data["First"])
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	data := PrepareTemplateData("ERROR 404", "code-404")
	RenderTemplateFiles(w, "code_404.html", data)
}

func handleBadRequest(w http.ResponseWriter, errorMessage string) {
	w.WriteHeader(http.StatusBadRequest)
	data := PrepareTemplateData("ERROR 400", "code-400")
	data["Error"] = errorMessage
	RenderTemplateFiles(w, "code_400.html", data)
}

func handleServerError(w http.ResponseWriter, errorMessage string) {
	w.WriteHeader(http.StatusInternalServerError)
	data := PrepareTemplateData("ERROR 500", "code-500")
	data["Error"] = errorMessage
	RenderTemplateFiles(w, "code_500.html", data)
}

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	// Check if the request is from a browser
	log.Printf("HandleDownload: Accept Header: %s\n", r.Header.Get("Accept"))
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		NotFoundHandler(w, r) // Redirect to custom 404 page
		return
	}

	// Get the current working directory
	rootDir, err := os.Getwd()
	if err != nil {
		http.Error(w, "Error getting current working directory", http.StatusInternalServerError)
		return
	}

	// Path to the file
	filePath := filepath.Join(rootDir, "temp", "output.txt")

	// Get file info to calculate Content-Length
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error checking file", http.StatusInternalServerError)
		return
	}

	// Get the size of the file
	fileSize := fileInfo.Size()

	// Set headers for download
	w.Header().Set("Content-Disposition", "attachment; filename=output.txt")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10)) // Convert file size to string

	// Serve the file
	http.ServeFile(w, r, filePath)
}

func HandleCheckFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleCheckFile: Accept Header: %s\n", r.Header.Get("Accept"))
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		NotFoundHandler(w, r) // Redirect to custom 404 page
		return
	}

	// Get the current working directory
	rootDir, err := os.Getwd()
	if err != nil {
		http.Error(w, "Error getting current working directory", http.StatusInternalServerError)
		return
	}

	// Path to the file
	filePath := filepath.Join(rootDir, "temp", "output.txt")

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// File exists
	w.WriteHeader(http.StatusOK)
}

func HandleExportZip(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleExportZip: Accept Header: %s\n", r.Header.Get("Accept"))
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		NotFoundHandler(w, r) // Redirect to custom 404 page
		return
	}
	// Get the current working directory
	rootDir, err := os.Getwd()
	if err != nil {
		http.Error(w, "Error getting current working directory", http.StatusInternalServerError)
		return
	}

	// Path to the file to be zipped
	filePath := filepath.Join(rootDir, "temp", "output.txt")

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Create an in-memory buffer to write the zip archive
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Add the file to the zip archive
	err = addFileToZip(zipWriter, filePath)
	if err != nil {
		http.Error(w, "Error creating zip file", http.StatusInternalServerError)
		return
	}

	// Close the zip writer
	err = zipWriter.Close()
	if err != nil {
		http.Error(w, "Error finalizing zip file", http.StatusInternalServerError)
		return
	}

	// Set headers for the zip file download
	w.Header().Set("Content-Disposition", `attachment; filename="output.zip"`)
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Length", strconv.Itoa(len(buf.Bytes())))

	// Write the zip file to the response
	_, err = w.Write(buf.Bytes())
	if err != nil {
		http.Error(w, "Error sending zip file", http.StatusInternalServerError)
		return
	}
}

func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	// Open the file to be added
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	// Create a zip file header based on the file info
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	header.Name = filepath.Base(filePath) // Use the file's base name for the zip entry
	header.Method = zip.Deflate           // Use deflate compression

	// Create a writer for the zip file entry
	zipFileWriter, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Copy the file's contents to the zip file
	_, err = io.Copy(zipFileWriter, file)
	if err != nil {
		return err
	}

	return nil
}
