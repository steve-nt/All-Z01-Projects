package backend

import (
	"fmt"
	"log"
	"os"
)

// logChan is a buffered channel used for logging messages asynchronously
var logChan = make(chan string, 100)


// init function starts a background goroutine to log messages to a file
func init() {
	go func() {
		// Open or create the history log file
		logFile, err := os.OpenFile("history.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			// Print error message if the log file cannot be opened
			fmt.Fprintln(os.Stderr, "Failed to open log file:", err)
			return
		}
		defer logFile.Close()

		// Create a logger instance for writing to the log file
		logger := log.New(logFile, "", log.LstdFlags)
		for msg := range logChan {
			logger.Println(msg)
		}
	}()
}


// LogHistory sends log messages to the log channel for asynchronous writing
func LogHistory(message string) {
	logChan <- message
}
