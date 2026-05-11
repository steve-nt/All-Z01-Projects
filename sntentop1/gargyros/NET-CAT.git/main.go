package main

import (
	"fmt"
	"os"
)

func main() {
	var port string

	if len(os.Args) == 1 {
		port = "8989"
	} else if len(os.Args) == 2 {
		port = os.Args[1]
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	server := NewServer(port)
	if err := server.Start(); err != nil {
		fmt.Println("Error:", err)
	}
}
