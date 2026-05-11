package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

const (
	defaultPort = 8989
	maxClients  = 10
)

type Server struct {
	clients    map[net.Conn]*Client
	messages   []Message
	mu         sync.RWMutex
	clientChan chan *Client
}

type Client struct {
	conn     net.Conn
	name     string
	messages chan Message
}

type Message struct {
	timestamp string
	username  string
	content   string
}

func main() {
	port := defaultPort

	if len(os.Args) > 1 {
		if len(os.Args) > 2 {
			fmt.Println("[USAGE]: ./TCPChat $port")
			os.Exit(1)
		}
		p, err := strconv.Atoi(os.Args[1])
		if err != nil || p < 1 || p > 65535 {
			fmt.Println("[USAGE]: ./TCPChat $port")
			os.Exit(1)
		}
		port = p
	}

	server := &Server{
		clients:    make(map[net.Conn]*Client),
		messages:   make([]Message, 0),
		clientChan: make(chan *Client),
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Listening on the port :%d\n", port)

	go server.handleMessages()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		if len(server.clients) >= maxClients {
			conn.Write([]byte("Chat is full. Please try again later.\n"))
			conn.Close()
			continue
		}

		go server.handleNewConnection(conn)
	}
} 