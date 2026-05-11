package main

import (
	"fmt"
	"net"
	"netcat/internal/server"
	"netcat/internal/types"
	"os"
)

func main() {
	args := os.Args[1:]
	port := "8989"

	if len(args) == 1 {
		port = args[0]
	} else if len(args) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	// Struct-based state logic
	state := &types.ChatState{
		Clients:     make(map[net.Conn]string),
		MessageHist: make([]string, 0),
	}

	server.StartServer(":"+port, state)
}
