# TCPChat

A concurrent TCP group chat server written in Go. Supports up to 10 simultaneous clients, chat history replay, and an optional terminal UI client — all over raw TCP, connectable with `nc`.

---

## Run

```bash
./build.sh          # build the binary

./TCPChat           # listen on default port 8989
./TCPChat 2525      # custom port

nc localhost 8989   # connect as a plain client
./TCPChat -tui      # connect with the TUI client
```

---

## Features

- Up to 10 concurrent connections; new connections beyond the limit are rejected cleanly
- Unique username enforced at join time — server prompts until a free name is entered
- Full chat history replayed to each new client on connect
- `/name NewUsername` command for runtime username changes
- All events (join, leave, rename, messages) logged to `server.log`
- Optional TUI client (`-tui`) built with bubbletea

---

## Message format

```
[2006-01-02 15:04:05][Username]: message text
```

System events (join/leave/rename) broadcast to all connected clients.

---

## Implementation notes

- `sync.RWMutex` protects the shared clients map and message history — read lock for broadcasts, write lock for mutations
- Each client runs in its own goroutine; disconnects are signaled via a channel to a dedicated `MonitorDisconnections` goroutine, keeping cleanup out of the per-client path
- Chat history is replayed before the client is registered in the map, so there's no race between history and live messages

---

## Project structure

```
net-cat/
├── main.go           # entry point, flag parsing, server/TUI dispatch
├── tui.go            # terminal UI client
├── server/
│   ├── server.go     # Server struct, init, logging
│   ├── connections.go # AcceptClient, MonitorDisconnections
│   ├── handlers.go   # input reading, command handling
│   ├── message.go    # message types and formatting
│   ├── client.go     # Client struct
│   └── utils.go
└── build.sh
```

---

## License

MIT
