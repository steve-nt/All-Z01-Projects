package main

import (
	"groupie-tracker/bin"
	"groupie-tracker/core"
	"log"
)

// init function is called before the main function
func init() {
	log.Println("Initializing application...")

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Initialize global variables if necessary
	if core.ShouldRecreateCacheFiles() {
		go bin.PopulateUniqueLocations()
	} else {
		core.InitializeGlobalVariables()
	}
}

func main() {
	bin.StartLRUEviction()
	router := bin.RegisterRoutes()

	server := core.SetupServer(router)
	core.StartServer(server)
}
