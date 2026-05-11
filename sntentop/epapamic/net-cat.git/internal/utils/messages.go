package utils

import (
	"net"
	"netcat/internal/types"
)

func Broadcast(message string, exclude string, state *types.ChatState) {
	for conn := range state.Clients {
		if exclude != "" && conn.RemoteAddr().String() == exclude {
			continue // skip the sender if needed
		}
		conn.Write([]byte(message + "\n"))
	}
}

func SendHistory(conn net.Conn, state *types.ChatState) {
	for _, msg := range state.MessageHist {
		conn.Write([]byte(msg + "\n"))
	}
}
