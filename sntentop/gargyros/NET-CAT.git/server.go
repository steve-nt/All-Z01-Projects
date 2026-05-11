package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// Welcome message displayed to every new client
const welcomeMessage = "Welcome to TCP-Chat!\n" +
	"         _nnnn_\n" +
	"        dGGGGMMb\n" +
	"       @p~qp~~qMb\n" +
	"       M|@||@) M|\n" +
	"       @,----.JM|\n" +
	"      JS^\\__/  qKL\n" +
	"     dZP        qKRb\n" +
	"    dZP          qKKb\n" +
	"   fZP            SMMb\n" +
	"   HZM            MMMM\n" +
	"   FqM            MMMM\n" +
	" __| \".        |\\dS\"qML\n" +
	" |    `.       | `' \\Zq\n" +
	"_)      \\.___.,|     .'\n" +
	"\\____   )MMMMMP|   .'\n" +
	"     `-'       `--'\n" +
	"[ENTER YOUR NAME]: "

// Client represents a connected user
type Client struct {
	conn    net.Conn
	name    string
	writeMu sync.Mutex
}

// write safely sends a message to the client
func (c *Client) write(text string) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	c.conn.Write([]byte(text))
}

// prompt returns the formatted prompt with timestamp and username
func (c *Client) prompt() string {
	now := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s][%s]:", now, c.name)
}

// Server represents the chat server
type Server struct {
	port    string
	clients map[*Client]bool
	history []string
	mu      sync.Mutex
}

// NewServer initializes a new server instance
func NewServer(port string) *Server {
	return &Server{
		port:    port,
		clients: make(map[*Client]bool),
		history: []string{},
	}
}

// Start begins listening for incoming TCP connections
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Println("Listening on the port :" + s.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go s.handleConnection(conn)
	}
}

// handleConnection manages communication with a single client
func (s *Server) handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)

	conn.Write([]byte(welcomeMessage))

	name, err := s.readName(reader, conn)
	if err != nil {
		conn.Close()
		return
	}

	client := &Client{
		conn: conn,
		name: name,
	}

	if err := s.addClient(client); err != nil {
		conn.Write([]byte(err.Error() + "\n"))
		conn.Close()
		return
	}

	s.sendHistory(client)
	s.broadcastSystem(client.name+" has joined our chat...", client)
	client.write(client.prompt())

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			s.removeClient(client)
			return
		}

		message = strings.TrimSpace(message)
		if message == "" {
			client.write(client.prompt())
			continue
		}

		formatted := s.formatMessage(client.name, message)
		s.saveHistory(formatted)

		s.broadcastMessage(formatted, client)
		client.write(client.prompt())
	}
}

// readName ensures the client provides a non-empty username
func (s *Server) readName(reader *bufio.Reader, conn net.Conn) (string, error) {
	for {
		name, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		name = strings.TrimSpace(name)
		if name == "" {
			conn.Write([]byte("[ENTER YOUR NAME]: "))
			continue
		}

		return name, nil
	}
}

// addClient adds a new client to the server safely
func (s *Server) addClient(client *Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.clients) >= 10 {
		return fmt.Errorf("Chat is full. Try again later.")
	}

	for existing := range s.clients {
		if existing.name == client.name {
			return fmt.Errorf("This name is already taken.")
		}
	}

	s.clients[client] = true
	return nil
}

func (s *Server) removeClient(client *Client) {
	s.mu.Lock()
	_, exists := s.clients[client]
	if exists {
		delete(s.clients, client)
	}
	s.mu.Unlock()

	if exists {
		s.broadcastSystem(client.name+" has left our chat...", client)
	}

	client.conn.Close()
}

func (s *Server) saveHistory(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = append(s.history, message)
}

func (s *Server) sendHistory(client *Client) {
	s.mu.Lock()
	historyCopy := make([]string, len(s.history))
	copy(historyCopy, s.history)
	s.mu.Unlock()

	for _, message := range historyCopy {
		client.write(message + "\n")
	}
}

func (s *Server) getClientsExcept(sender *Client) []*Client {
	s.mu.Lock()
	defer s.mu.Unlock()

	clients := []*Client{}

	for client := range s.clients {
		if client != sender {
			clients = append(clients, client)
		}
	}

	return clients
}

func (s *Server) broadcastMessage(message string, sender *Client) {
	clients := s.getClientsExcept(sender)

	for _, client := range clients {
		client.write("\n" + message + "\n" + client.prompt())
	}
}

func (s *Server) broadcastSystem(message string, sender *Client) {
	clients := s.getClientsExcept(sender)

	for _, client := range clients {
		client.write("\n" + message + "\n" + client.prompt())
	}
}

func (s *Server) formatMessage(name, message string) string {
	now := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s][%s]:%s", now, name, message)
}
