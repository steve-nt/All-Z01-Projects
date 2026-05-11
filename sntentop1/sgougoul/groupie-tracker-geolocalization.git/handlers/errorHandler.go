package handlers

import (
	"log"
	"net/http"
)

type HTTPError struct {
	Status  int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}

type Apphandler func(w http.ResponseWriter, r *http.Request) error

func (fn Apphandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Printf("http error : %v", err)
		if httpErr, ok := err.(*HTTPError); ok {
			ErrorHandler(w, r, httpErr.Status, httpErr.Message)
		} else {
			ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		}
	}
}

// ErrorHandler handles different HTTP errors and displays appropriate messages
func ErrorHandler(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

}

// Handle 404 errors(Page Not Found)
func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	ErrorHandler(w, r, http.StatusNotFound, "Page not found.")
}
