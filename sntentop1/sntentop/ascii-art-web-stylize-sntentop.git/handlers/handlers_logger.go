package handlers

import (
	"fmt"
	"time"
)

var LogHistory []string

func LogEventsRecord(eventType, message string) { //Record all logs
	entry := fmt.Sprintf("[%s] %s: %s", time.Now().Format("2006-01-02 15:04:05"), eventType, message)
	LogHistory = append(LogHistory, entry)
}

func LogEventsPrint() { // Print the full history
	fmt.Println("\n===== Log History =====")
	for _, h := range LogHistory {
		fmt.Println(h)
	}
}
