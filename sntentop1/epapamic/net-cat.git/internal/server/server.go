package server

import (
	"fmt"
	"log"
	"net"
	"netcat/internal/client"
	"netcat/internal/types"
	"netcat/internal/utils"
)

func StartServer(port string, state *types.ChatState) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("Listening on port :", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		state.Mutex.Lock()
		if utils.MaxConnCheck(conn, state) {
			conn.Close()
			state.Mutex.Unlock()
			continue
		}
		state.Mutex.Unlock()

		go client.HandleConnection(conn, state)
	}
}
