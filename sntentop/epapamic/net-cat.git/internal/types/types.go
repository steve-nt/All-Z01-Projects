package types

import (
	"net"
	"sync"
)

type ChatState struct {
	Clients     map[net.Conn]string // Shared across client/utils
	Mutex       sync.Mutex          // Shared lock
	MessageHist []string            // Shared message history
}
