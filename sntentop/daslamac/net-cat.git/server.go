package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorPink   = "\033[38;5;219m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

const welcomeLogo = `Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    '.       | '  \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     '-'       '--'
`

func init() {
	// Set up logging to file
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	log.SetOutput(logFile)
}

func (s *Server) handleNewConnection(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte(welcomeLogo))
	conn.Write([]byte(colorPink + "[ENTER YOUR NAME]: " + colorReset))

	reader := bufio.NewReader(conn)
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading name: %v", err)
		return
	}

	name = strings.TrimSpace(name)
	if name == "" {
		conn.Write([]byte("Name cannot be empty\n"))
		return
	}

	// Check if name is already taken
	s.mu.RLock()
	for _, client := range s.clients {
		if client.name == name {
			s.mu.RUnlock()
			conn.Write([]byte("Name already taken. Please try another name.\n"))
			conn.Write([]byte(colorPink + "[ENTER YOUR NAME]: " + colorReset))
			newName, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error reading new name: %v", err)
				return
			}
			newName = strings.TrimSpace(newName)
			if newName == "" {
				conn.Write([]byte("Name cannot be empty\n"))
				return
			}
			name = newName
			break
		}
	}
	s.mu.RUnlock()

	client := &Client{
		conn:     conn,
		name:     name,
		messages: make(chan Message, 10),
	}

	// Send previous messages to new client
	s.mu.RLock()
	for _, msg := range s.messages {
		client.conn.Write([]byte(formatMessage(msg)))
	}
	s.mu.RUnlock()

	// Add client to server
	s.mu.Lock()
	s.clients[conn] = client
	s.mu.Unlock()

	// Broadcast join message
	joinMsg := Message{
		timestamp: time.Now().Format("2006-01-02 15:04:05"),
		content:   fmt.Sprintf("%s has joined our chat...\n", client.name),
	}
	s.broadcast(joinMsg)
	log.Printf("Client connected: %s from %s", client.name, conn.RemoteAddr())

	// Start reading messages from client
	go s.readClientMessages(client)

	// Block until client disconnects
	<-client.messages
}

func (s *Server) handleMessages() {
	for client := range s.clientChan {
		s.mu.Lock()
		delete(s.clients, client.conn)
		s.mu.Unlock()

		leaveMsg := Message{
			timestamp: time.Now().Format("2006-01-02 15:04:05"),
			content:   fmt.Sprintf("%s has left our chat...\n", client.name),
		}
		s.broadcast(leaveMsg)
		log.Printf("Client disconnected: %s", client.name)
	}
}

func (s *Server) readClientMessages(client *Client) {
	reader := bufio.NewReader(client.conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading message from %s: %v", client.name, err)
			break
		}

		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		// Check for name change command
		if strings.HasPrefix(message, "/name ") {
			newName := strings.TrimSpace(strings.TrimPrefix(message, "/name "))
			if newName == "" {
				client.conn.Write([]byte("Name cannot be empty\n"))
				continue
			}

			// Check if name is already taken
			nameExists := false
			s.mu.RLock()
			for _, c := range s.clients {
				if c.name == newName && c != client {
					nameExists = true
					break
				}
			}
			s.mu.RUnlock()

			if nameExists {
				client.conn.Write([]byte("Name already taken. Please try another name.\n"))
				continue
			}

			oldName := client.name
			s.mu.Lock()
			client.name = newName
			s.mu.Unlock()

			nameChangeMsg := Message{
				timestamp: time.Now().Format("2006-01-02 15:04:05"),
				content:   fmt.Sprintf("%s changed name to %s\n", oldName, newName),
			}
			s.broadcast(nameChangeMsg)
			log.Printf("Client %s changed name to %s", oldName, newName)
			continue
		}

		msg := Message{
			timestamp: time.Now().Format("2006-01-02 15:04:05"),
			username:  client.name,
			content:   message,
		}

		s.mu.Lock()
		s.messages = append(s.messages, msg)
		s.mu.Unlock()

		s.broadcast(msg)
		log.Printf("Message from %s: %s", client.name, message)
	}

	close(client.messages)
	s.clientChan <- client
}

func (s *Server) broadcast(msg Message) {
	formattedMsg := formatMessage(msg)

	s.mu.RLock()
	for _, client := range s.clients {
		client.conn.Write([]byte(formattedMsg))
	}
	s.mu.RUnlock()
}

func formatMessage(msg Message) string {
	if msg.username == "" {
		return colorYellow + msg.content + colorReset
	}
	return fmt.Sprintf(colorGreen + "[%s][%s]" + colorReset + ":%s\n", msg.timestamp, msg.username, msg.content)
} 