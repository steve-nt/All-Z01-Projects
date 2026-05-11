package utils

import (
	"net"
	"netcat/internal/types"
)

func MaxConnCheck(conn net.Conn, state *types.ChatState) bool {
	if len(state.Clients) >= 10 {
		conn.Write([]byte("Server is full. Try again later.\n"))
		return true
	}
	return false
}
