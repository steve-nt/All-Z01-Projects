package handlers

import (
	"archive/zip"   // Provides functionality for creating and extracting ZIP archives.
	"bytes"         // Implements functions for byte manipulation.
	"io"            // Contains basic interfaces for I/O primitives.
	"log"           // Used for logging messages.
	"net/http"      // Provides HTTP client and server implementations.
	"os"            // Offers functions to perform OS-level operations like file and directory management.
	"path/filepath" // Manages file path manipulations in a cross-platform way.
	"strconv"       // Provides utilities for converting strings to numbers and vice versa.
	"strings"       // Implements string manipulation functions.
)

// PgHome handles requests for the home page.
func PgHome(w http.ResponseWriter, r *http.Request) {
	data := PrepareTemplateData("HOME", "home-page")                  // Calls a helper to prepare data for the "home-page" template.
	log.Println("requesting_page: HOME")                              // Logs that the "HOME" page is being requested.
	RenderTemplateFiles(w, "PgHome.html", data)                       // Renders the HTML template for the home page.
	log.Printf("Rendering template %v successfully\n", data["Title"]) // Logs successful rendering of the template.
}

func PgProject(w http.ResponseWriter, r *http.Request) {
	data := PrepareTemplateData("PROJECT", "project-page")            // Prepares data for the "project-page" template.
	log.Println("requesting_page: PROJECT")                           // Logs the page request.
	RenderTemplateFiles(w, "PgProject.html", data)                    // Renders the HTML template for the home page.
	log.Printf("Rendering template %v successfully\n", data["Title"]) // Logs successful rendering of the template.
}

func PgTeam(w http.ResponseWriter, r *http.Request) {
	data := PrepareTemplateData("TEAM", "team-page")                  // Prepares data for the "team-page" template.
	log.Println("requesting_page: TEAM")                              // Logs the page request.
	RenderTemplateFiles(w, "PgTeam.html", data)                       // Renders the template for the team page.
	log.Printf("Rendering template %v successfully\n", data["Title"]) // Logs successful rendering of the template.
}

func PgConverter(w http.ResponseWriter, r *http.Request) {
	data := PrepareTemplateData("CONVERTER", "converter-page") // Prepares data for the "converter-page"
	if r.Method == http.MethodPost {                           // Checks if the request method is POST.

		text := r.FormValue("ascii-data") // Retrieves the ASCII data from the form.
		banner := r.FormValue("banner")   // Retrieves the banner style from the form.

		final, renderStatus := inputValidation(text) // Validates and processes the input data.
		log.Println("Render Status:", renderStatus, "Rendered Text:", final)
		if renderStatus == 400 { // Handles bad input.
			log.Println("400 Error", "Invalid Input: "+text)
			handleBadRequest(w, "Invalid Input: "+text) // Renders a "400 Bad Request" error page.
			return
		}

		if renderStatus != 200 { // Handles other errors.
			data["Error"] = final
		} else { // Successful input validation.
			output, statusCode := asciiArt(final, banner) // Generates ASCII art.
			if statusCode == 500 {                        // Handles server errors during ASCII art generation.
				log.Println("500 Error", "Failed to read font file for font: "+banner)
				handleServerError(w, "Failed to generate ASCII art") // Renders a "500 Internal Server Error" page.
				return
			}

			if statusCode != 200 { // Handles unexpected errors.
				data["Error"] = "Failed to generate ASCII art. Please check your input and try again."
			} else { // Success case.
				data["First"] = output             // Stores the ASCII art result.
				fpath, err := ExportToFile(output) // Exports the result to a file.
				if err != nil {                    // Logs errors during file export.
					log.Printf("Error: %v", err)
				} else { // Logs successful file export.
					log.Printf("File successfully exported to: %s", fpath)
				}

			}
		}
	}
	log.Println("requesting-page: CONVERTER")                         // na ginei me LogEventsRecord???  (Logs the page request.)
	RenderTemplateFiles(w, "PgConverter.html", data)                  // Renders the converter page template.
	log.Printf("Rendering template %v successfully\n", data["Title"]) // Logs successful rendering.
	log.Printf("The output:\n %v\n", data["First"])                   // Logs the generated ASCII art.
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound) // Sends an HTTP 404 status code to indicate "Not Found".
	// Prepares the data needed to render the error page.
	// "ERROR 404" is the page title, and "code-404" is the template name.
	data := PrepareTemplateData("ERROR 404", "code-404")
	// Renders the error page using the specified template and data.
	RenderTemplateFiles(w, "code_404.html", data)
}

func handleBadRequest(w http.ResponseWriter, errorMessage string) {
	// Sends an HTTP 400 status code to indicate a "Bad Request".
	w.WriteHeader(http.StatusBadRequest)
	// Prepares data for the 400 error page, including an error message.
	data := PrepareTemplateData("ERROR 400", "code-400")
	data["Error"] = errorMessage
	// Renders the error page with the provided template and data.
	RenderTemplateFiles(w, "code_400.html", data)
}

func handleServerError(w http.ResponseWriter, errorMessage string) {
	// Sends an HTTP 500 status code for "Internal Server Error".
	w.WriteHeader(http.StatusInternalServerError)
	// Prepares data for the 500 error page, including the error message.
	data := PrepareTemplateData("ERROR 500", "code-500")
	data["Error"] = errorMessage // Inserts the error message into the data.
	// Renders the error page using the specified template and data.
	RenderTemplateFiles(w, "code_500.html", data)
}

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	// Logs the Accept header from the request for debugging purposes.
	log.Printf("HandleDownload: Accept Header: %s\n", r.Header.Get("Accept"))
	// Checks if the request is from a browser. If so, calls `NotFoundHandler`.
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		NotFoundHandler(w, r) // Redirect to custom 404 page
		return
	}

	// Gets the current working directory.
	rootDir, err := os.Getwd()
	if err != nil {
		http.Error(w, "Error getting current working directory", http.StatusInternalServerError)
		return
	}

	// Constructs the path to the file to be downloaded
	filePath := filepath.Join(rootDir, "temp", "output.txt")

	// Checks if the file exists. If not, returns a 404 error.
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error checking file", http.StatusInternalServerError)
		return
	}

	// Retrieves the file's size to include in the response headers.
	fileSize := fileInfo.Size()

	// Sets response headers to facilitate file download.
	w.Header().Set("Content-Disposition", "attachment; filename=output.txt")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10)) // Convert file size to string

	// Serves the file to the client.
	http.ServeFile(w, r, filePath)
}

func HandleCheckFile(w http.ResponseWriter, r *http.Request) {
	// Logs the Accept header for debugging purposes.
	log.Printf("HandleCheckFile: Accept Header: %s\n", r.Header.Get("Accept"))
	// Redirects browser requests to a custom 404 page.
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

	// Constructs the path to the file to check.
	filePath := filepath.Join(rootDir, "temp", "output.txt")

	// Checks if the file exists.
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// If the file exists, sends a 200 status code.
	w.WriteHeader(http.StatusOK)
}

func HandleExportZip(w http.ResponseWriter, r *http.Request) {
	// Logs the Accept header for debugging purposes.
	log.Printf("HandleExportZip: Accept Header: %s\n", r.Header.Get("Accept"))
	// Redirects browser requests to a custom 404 page.
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

	// Constructs the path to the file to be zipped.
	filePath := filepath.Join(rootDir, "temp", "output.txt")

	// Checks if the file exists. If not, sends a 404 error.
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Creates an in-memory buffer to write the ZIP archive.
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
	// Opens the file to be added to the ZIP archive.
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Retrieves file information to create a ZIP header.
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	// Creates a ZIP file header using the file information.
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	header.Name = filepath.Base(filePath) // Sets the file name in the ZIP archive.
	header.Method = zip.Deflate           // Compresses the file using the DEFLATE algorithm.

	// Creates a writer for the ZIP file entry.
	zipFileWriter, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Copies the file's contents to the ZIP archive.
	_, err = io.Copy(zipFileWriter, file)
	if err != nil {
		return err
	}

	return nil
}
