package client

import (
	"bufio"
	"fmt"
	"net"
	"netcat/internal/types"
	"netcat/internal/utils"
	"strings"
	"time"
)

func HandleConnection(conn net.Conn, state *types.ChatState) {
	defer conn.Close()

	// Send banner and prompt for name
	conn.Write([]byte(utils.LoadBanner() + "\n[ENTER YOUR NAME]: "))

	reader := bufio.NewReader(conn)
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read name:", err)
		return
	}
	name = strings.TrimSpace(name)

	if name == "" {
		conn.Write([]byte("Name cannot be empty.\n"))
		return
	}

	// Register the client
	state.Mutex.Lock()
	if utils.MaxConnCheck(conn, state) {
		state.Mutex.Unlock()
		return
	}
	state.Clients[conn] = name
	utils.SendHistory(conn, state)
	utils.Broadcast(fmt.Sprintf("%s has joined our chat...", name), "", state)
	state.Mutex.Unlock()

	// Let the client know they're ready to chat
	conn.Write([]byte("\n--- Chat Ready ---\nYou can now start typing your message:\n"))

	// Read incoming messages from the client
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		if strings.TrimSpace(msg) == "" {
			continue
		}
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		formatted := fmt.Sprintf("[%s][%s]:%s", timestamp, name, msg)

		state.Mutex.Lock()
		state.MessageHist = append(state.MessageHist, formatted)
		utils.Broadcast(formatted, conn.LocalAddr().String(), state)
		state.Mutex.Unlock()
	}

	// Handle disconnect
	state.Mutex.Lock()
	utils.Broadcast(fmt.Sprintf("%s has left our chat...", name), "", state)
	delete(state.Clients, conn)
	state.Mutex.Unlock()
}
