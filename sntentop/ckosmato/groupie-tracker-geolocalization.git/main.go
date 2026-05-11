package main

import (
	"fmt"
	"gtracker/router"
	"net/http"
)

func main() {
	router.InitRoutes()
	fmt.Println("Server is running on localhost:8080")
	http.ListenAndServe(":8080", nil)
}
