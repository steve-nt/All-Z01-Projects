# TCPChat

TCPChat is a NetCat-like group chat application written in Go.

The project uses a server-client architecture. The server listens on a TCP port and multiple clients can connect to it using `nc`.

---

## Features

- TCP connection between one server and multiple clients
- Maximum 10 simultaneous clients
- Client name requirement
- Unique usernames
- Clients can send messages to the group chat
- Empty messages are ignored
- Messages include timestamp and username
- New clients receive the previous chat history
- Clients are notified when another client joins
- Clients are notified when another client leaves
- If one client disconnects, the rest stay connected
- Uses goroutines
- Uses mutexes for safe concurrency
- Default port is `8989`

---

## Message Format

    [YYYY-MM-DD HH:MM:SS][username]:message

Example:

    [2026-04-29 11:45:10][John]:hello

---

## Usage

Build:

    go build -o TCPChat .

Run with default port:

    ./TCPChat

Run with custom port:

    ./TCPChat 2525

Wrong usage:

    ./TCPChat 2525 localhost

Output:

    [USAGE]: ./TCPChat $port

---

## Connect as a Client

    nc localhost 2525

Then enter your name:

    [ENTER YOUR NAME]: John

---

## Example

Terminal 1:

    ./TCPChat 2525

Terminal 2:

    nc localhost 2525

Terminal 3:

    nc localhost 2525

---

## Project Structure

    .
    ├── main.go
    ├── server.go
    └── README.md

---

## Packages Used

- bufio
- fmt
- net
- os
- strings
- sync
- time

---

## Goroutines

The project uses goroutines to handle multiple clients:

    go s.handleConnection(conn)

---

## Mutexes

The project uses mutexes to protect shared data:

- clients list
- chat history
- client writes

---

## How It Works

1. Server starts and listens on a TCP port
2. Clients connect using `nc`
3. Server asks for username
4. Client joins chat
5. History is sent
6. Messages are broadcast
7. Clients are notified on join/leave

---

## Author

Created as part of a Go networking project.