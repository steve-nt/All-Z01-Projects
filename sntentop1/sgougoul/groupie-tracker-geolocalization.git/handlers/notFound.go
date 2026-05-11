package handlers

import (
	"net/http"
)

func CustomNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>404 - Page Not Found</title>
            <link href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" rel="stylesheet">
            <link rel="stylesheet" href="/static/home.css">
        </head>
        <body>
            <div class="container">
                <h1 class="my-4 text-center">404 - Page Not Found</h1>
                <p class="text-center">The page you requested doesn't exist.</p>
                <div class="text-center">
                    <a href="/" class="btn btn-primary">Go Home</a>
                </div>
            </div>
        </body>
        </html>
    `))
}
