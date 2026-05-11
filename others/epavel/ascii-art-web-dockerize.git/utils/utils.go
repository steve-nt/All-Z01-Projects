package utils

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Function to render the error page
func RenderErrorPage(w http.ResponseWriter, errorMessage string) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", "error.html"))
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		log.Println("Error loading template:", err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, errorMessage)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error rendering error template:", err)
	}
}

func ParseBanner(banner string) ([]string, error) {
	_, current_file, _, _ := runtime.Caller(0)
	current_dir := filepath.Dir(current_file)
	mode := filepath.Join(current_dir, "banners", banner+".txt")
	reference, err := os.ReadFile(mode)
	if err != nil {
		return nil, errors.New("error reading file")
	}
	lines := strings.Split(string(reference), "\n")
	return lines, nil
}

func Indexing(input string) ([]int, error) {
	indeces := []int{}
	found_new_line := false
	for i := 0; i < len(input); i++ {
		if found_new_line {
			found_new_line = false
			continue
		}
		byte_value := int(input[i])
		if byte_value < 32 || byte_value > 126 {
			return nil, errors.New("invalid character. Allowed characters are between abc ... z and ABC ... Z and special characters")
		}
		if string(input[i]) == "\\" && i == len(input)-1 {
			indeces = append(indeces, (byte_value-32)*9+1)
			continue
		}
		if string(input[i]) == "\\" && string(input[i+1]) == "n" {
			found_new_line = true
			indeces = append(indeces, -1) // assigning -1 to \n as a reference for future handling
			continue
		}
		//	compensation for line between every 8-line ASCII char(standard.txt)
		indeces = append(indeces, (byte_value-32)*9+1)
	}
	return indeces, nil
}

func NewLineHandling(indeces []int) [][]int {
	// Create a slice of slices to hold sub-slices
	var subSlices [][]int
	var currentSlice []int

	for _, index := range indeces {
		if index == -1 {
			// Append the current slice to subSlices and reset it
			subSlices = append(subSlices, currentSlice)
			currentSlice = []int{}
		} else {
			// Add the index to the current slice
			currentSlice = append(currentSlice, index)
		}
	}

	// Append the last slice if it's not empty
	if len(currentSlice) > 0 {
		subSlices = append(subSlices, currentSlice)
	}
	return subSlices
}

func OutputAscii(indeces []int, lines []string) string {
	output := ""
	if len(indeces) == 0 {
		output += "\n"
		return output
	}
	for counter := 0; counter < 8; counter++ {
		for i := 0; i < len(indeces); i++ {
			output += lines[indeces[i]+counter]
		}
		if counter < 8 {
			output += "\n"
		}
	}
	return output
}
