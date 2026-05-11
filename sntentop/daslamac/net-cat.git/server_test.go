package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func TestServerBasicFunctionality(t *testing.T) {
	// Start server
	server := &Server{
		clients:    make(map[net.Conn]*Client),
		messages:   make([]Message, 0),
		clientChan: make(chan *Client),
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go server.handleNewConnection(conn)
		}
	}()

	// Connect test client
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		t.Fatalf("Failed to connect test client: %v", err)
	}
	defer conn.Close()

	// Read welcome message
	reader := bufio.NewReader(conn)
	welcome, err := reader.ReadString(':')
	if err != nil {
		t.Fatalf("Failed to read welcome message: %v", err)
	}
	if !strings.Contains(welcome, "Welcome to TCP-Chat!") {
		t.Error("Welcome message not received")
	}

	// Send client name
	name := "TestUser\n"
	conn.Write([]byte(name))
	time.Sleep(100 * time.Millisecond)

	// Verify client was added
	server.mu.RLock()
	if len(server.clients) != 1 {
		t.Error("Client was not added to server")
	}
	server.mu.RUnlock()
}

func TestMaxClients(t *testing.T) {
	server := &Server{
		clients:    make(map[net.Conn]*Client),
		messages:   make([]Message, 0),
		clientChan: make(chan *Client),
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			if len(server.clients) >= maxClients {
				conn.Write([]byte("Chat is full. Please try again later.\n"))
				conn.Close()
				continue
			}
			go server.handleNewConnection(conn)
		}
	}()

	// Try to connect more than maxClients
	for i := 0; i <= maxClients; i++ {
		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			t.Fatalf("Failed to connect test client: %v", err)
		}
		defer conn.Close()

		reader := bufio.NewReader(conn)
		if i == maxClients {
			// Last connection should be rejected
			msg, _ := reader.ReadString('\n')
			if !strings.Contains(msg, "Chat is full") {
				t.Error("Max clients limit not enforced")
			}
		} else {
			// Send client name for valid connections
			reader.ReadString(':') // Read welcome message
			conn.Write([]byte(fmt.Sprintf("TestUser%d\n", i)))
		}
		time.Sleep(100 * time.Millisecond)
	}

	server.mu.RLock()
	if len(server.clients) > maxClients {
		t.Error("Server accepted more than maximum allowed clients")
	}
	server.mu.RUnlock()
} 