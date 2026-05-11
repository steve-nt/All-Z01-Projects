// Package handlers defines a collection of functionalities for handling specific tasks,
// in this case, logging events.
package handlers

// Importing external Go standard packages.
import (
	"fmt"  // The "fmt" package provides formatted I/O with functions for printing and string formatting.
	"time" // The "time" package provides functionality to manipulate and format time and dates.
)

// LogHistory is a global variable that stores a list of log entries.
// It is a slice of strings where each string represents a single log entry.
var LogHistory []string

// LogEventsRecord logs an event with a timestamp, event type, and message.
// The function accepts two parameters:
// - eventType: A string indicating the type of event (e.g., "INFO", "ERROR").
// - message: A string providing details about the event.
func LogEventsRecord(eventType, message string) { //Record all logs
	// time.Now() returns the current date and time.
	// .Format() formats the time in a specific layout defined by "2006-01-02 15:04:05".
	// The format string "2006-01-02 15:04:05" is the reference layout for time formatting in Go.
	entry := fmt.Sprintf("[%s] %s: %s", time.Now().Format("2006-01-02 15:04:05"), eventType, message)
	// fmt.Sprintf formats a string using placeholders. In this case:
	// [%s] - the timestamp (formatted current time),
	// %s - the event type, and
	// %s - the event message.
	// Append the formatted log entry to the LogHistory slice.
	LogHistory = append(LogHistory, entry)
}

// LogEventsPrint prints all recorded log entries to the console.
// It does not take any parameters and does not return a value.
func LogEventsPrint() { // Print the full history
	// fmt.Println prints a string to the console with a newline at the end.
	fmt.Println("\n===== Log History =====") // Print a header to separate logs visually.
	// Iterate over each entry in the LogHistory slice
	for _, h := range LogHistory {
		// Print the log entry to the console.
		fmt.Println(h)
	}
}
