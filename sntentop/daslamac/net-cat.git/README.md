# TCP Chat 

A Go implementation of a TCP chat server similar to NetCat, allowing multiple clients to connect and communicate in a group chat environment.

## Features

- TCP connection between server and multiple clients (1-to-many)
- Client name requirement with duplicate name prevention
- Maximum 10 concurrent connections
- Message broadcasting to all connected clients
- Empty message filtering
- Timestamped messages with sender identification (color-coded)
- Message history for new clients
- Join/Leave notifications
- Graceful client disconnection handling
- Name change functionality `/name <new-name>`
- Server logging to file
- Colored text interface

## Usage

### Building the Project

```bash
# Build the executable
./build.sh
```

### Starting the Server

```bash
# Default port (8989)
./TCPChat

# Custom port
./TCPChat 2525
```

### Connecting as a Client

Use the `nc` command or any TCP client to connect:

```bash
nc localhost 8989
```

### Using the Chat

- Enter your name when prompted
- Type messages and press Enter to send
- Change your name with `/name <new-name>`
- Exit by closing the terminal or pressing Ctrl+C

## Message Format

Messages are formatted as:
```
[YYYY-MM-DD HH:MM:SS][username]:message
```

## Implementation Details

- Uses Go's built-in `net` package for TCP networking
- Implements concurrent client handling using goroutines
- Uses channels for client message broadcasting
- Implements mutex locks for thread-safe operations
- Color-coded interface for better user experience:
  - Pink: Input prompts
  - Green: Timestamps and usernames
  - Yellow: System messages
- Server logging to "server.log"
- Includes unit tests for core functionality

## Project Structure

- `main.go`: Entry point and server initialization
- `server.go`: Server implementation and client handling
- `server_test.go`: Unit tests
- `build.sh`: Build script

## Testing

Run the tests using:

```bash
go test -v
```

## Requirements

- Go 1.x or higher

## Error Handling

- Graceful handling of client disconnections
- Input validation for client names
- Connection limit enforcement
- Error logging for debugging
- Duplicate name detection 