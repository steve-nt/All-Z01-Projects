# TCP Chat Server

A simple, concurrent TCP chat server built in Go that allows multiple clients to connect and chat in real-time.

## Features

- **Multi-client support**: Up to 10 concurrent connections
- **Real-time messaging**: Instant message broadcasting to all connected clients
- **Message history**: New clients receive chat history upon joining
- **Connection management**: Automatic cleanup when clients disconnect
- **ASCII art banner**: Welcoming banner displayed to new clients
- **Timestamped messages**: All messages include timestamp and sender name
- **Thread-safe**: Concurrent access handled with mutex locks

## Project Structure

```
netcat/
├── main.go
├── internal/
│   ├── server/
│   │   └── server.go
│   ├── client/
│   │   └── client.go
│   ├── types/
│   │   └── types.go
│   └── utils/
│       ├── banner.go
│       ├── connection.go
│       └── messaging.go
└── TCPChat (executable)
```

## Installation

### Prerequisites
- Go 1.16 or higher

### Building from Source
```bash
# Clone or download the project
# Navigate to the project directory
go mod tidy
go build -o TCPChat
```

## Usage

### Starting the Server
```bash
# Default port (8989)
./TCPChat

# Custom port
./TCPChat 3000
```

### Command Line Arguments
- **No arguments**: Server starts on port 8989
- **One argument**: Server starts on the specified port
- **Multiple arguments**: Shows usage message and exits

### Connecting to the Server
Use any TCP client like `telnet` or `nc`:

```bash
# Connect to default port
telnet localhost 8989

# Connect to custom port
telnet localhost 3000

# Using netcat
nc localhost 8989
```

## How It Works

1. **Server Startup**: The server binds to the specified port and starts listening for connections
2. **Client Connection**: When a client connects, they see a welcome banner and are prompted to enter their name
3. **Name Registration**: The client's name is registered and they receive the chat history
4. **Join Notification**: Other clients are notified when someone joins
5. **Message Broadcasting**: All messages are timestamped, stored in history, and broadcast to all clients
6. **Disconnect Handling**: When a client disconnects, others are notified and the connection is cleaned up

## Message Format

Messages are formatted as:
```
[YYYY-MM-DD HH:MM:SS][Username]:Message content
```

Example:
```
[2024-01-15 14:30:25][Alice]:Hello everyone!
```

## Server Limits

- **Maximum connections**: 10 concurrent clients
- **Connection overflow**: Additional clients receive "Server is full" message
- **Empty names**: Not allowed, connection will be terminated

## Architecture

### Core Components

- **`main.go`**: Entry point, argument parsing, and server initialization
- **`server/server.go`**: TCP server setup and connection acceptance
- **`client/client.go`**: Individual client connection handling
- **`types/types.go`**: Shared data structures and state management
- **`utils/`**: Utility functions for messaging, connections, and banner display

### State Management

The application uses a centralized `ChatState` struct with:
- `Clients`: Map of active connections to usernames
- `MessageHist`: Slice storing chat history
- `Mutex`: Synchronization for thread-safe operations

## Error Handling

- Invalid port numbers
- Network connection failures
- Client disconnections
- Maximum connection limits
- Empty username validation

## Example Session

```
$ ./TCPChat 8080
Listening on port : :8080

# Client 1 connects
Welcome to TCP-Chat!
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
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'

[ENTER YOUR NAME]: Alice

--- Chat Ready ---
You can now start typing your message:
Hello everyone!

# Client 2 connects and sees:
Alice has joined our chat...
[2024-01-15 14:30:25][Alice]:Hello everyone!
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is open source. Feel free to use and modify as needed.

## Troubleshooting

### Common Issues

**Port already in use**
```
Error: listen tcp :8989: bind: address already in use
```
Solution: Use a different port or kill the process using the port

**Permission denied**
```
Error: listen tcp :80: bind: permission denied
```
Solution: Use a port above 1024 or run with appropriate permissions

**Connection refused**
```
Error: dial tcp 127.0.0.1:8989: connect: connection refused
```
Solution: Ensure the server is running and the port is correct

### Testing the Server

```bash
# Terminal 1: Start server
./TCPChat 2525

# Terminal 2: Connect first client
telnet localhost 2525

# Terminal 3: Connect second client
telnet localhost 2525
```

Send messages from each client to test real-time broadcasting and message history functionality.