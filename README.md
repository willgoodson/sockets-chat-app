# Sockets Chat App

A terminal chat room in Go, built directly on TCP sockets. A central server relays every message to all connected clients.

## Overview

This was a learning project for Go's networking and concurrency primitives. It uses no HTTP or WebSocket library, only `net.Listen` / `net.Dial` and goroutines. Messages are serialized with `encoding/gob`.

- The **server** accepts TCP connections and starts a goroutine to read from each client. Every message goes into a shared channel, and a broadcaster goroutine sends it on to all connected clients.
- The **client** asks for a username, sends each line you type to the server, and prints messages from other users. Each client generates a UUID so it can skip its own messages when they're echoed back.

## Tech Stack

| Layer         | Technology                         |
|---------------|------------------------------------|
| Language      | Go                                 |
| Transport     | Raw TCP (`net` package)            |
| Serialization | `encoding/gob`                     |
| IDs           | [google/uuid](https://github.com/google/uuid) |

## Project Structure

```
server/server.go   Accepts connections and broadcasts messages
client/client.go   Interactive terminal client
```

## Getting Started

**Prerequisites:** Go 1.20+

```bash
# Terminal 1: start the server (listens on localhost:9988)
go run ./server

# Terminals 2..n: start one client per user
go run ./client
```

Type a username, then start chatting.

## Future Work

- [ ] Announce join/leave events to the room
- [ ] Make the host and port configurable with flags so it can run across machines
- [ ] Move the shared `message` type into its own package used by both server and client
- [ ] Add a `/quit` command and graceful shutdown
- [ ] Support multiple rooms and private messages
- [ ] Add a terminal UI (for example [Bubble Tea](https://github.com/charmbracelet/bubbletea)) so incoming messages don't overwrite the input prompt
- [ ] Add TLS
