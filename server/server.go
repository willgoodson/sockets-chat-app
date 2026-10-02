package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sync"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9988"
	SERVER_TYPE = "tcp"
)

type message struct {
	Id       string
	Username string
	Message  string
}

func main() {
	server()
}

// Set of connected clients, shared between the accept loop and the broadcaster
var (
	clients    = make(map[net.Conn]struct{})
	clients_mu sync.Mutex
)

func server() {
	msg_channel := make(chan message)

	fmt.Println("Server Running...")
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error Listening: ", err.Error())
		os.Exit(1)
	}
	defer server.Close()

	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
	fmt.Println("Waiting on client...")

	go broadcast_clients(msg_channel)

	for {
		// Wait for connections
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error Accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Client Connected\n")
		// Add client connection to map
		clients_mu.Lock()
		clients[connection] = struct{}{}
		clients_mu.Unlock()
		// Start routine to handle incoming client messages
		go read_client(connection, msg_channel)
	}
}

// Handle incoming messages from clients
func read_client(connection net.Conn, msg_channel chan message) {
	defer func() {
		clients_mu.Lock()
		delete(clients, connection)
		clients_mu.Unlock()
		connection.Close()
	}()
	decoder := gob.NewDecoder(connection)
	for {
		var incoming_msg message
		err := decoder.Decode(&incoming_msg)
		if err != nil {
			fmt.Println("Error Decoding: ", err.Error())
			break
		}
		fmt.Printf("%s> %s", incoming_msg.Username, incoming_msg.Message)
		msg_channel <- incoming_msg
	}
}

// Handle outgoing messages to clients
func broadcast_clients(msg_channel chan message) {
	encoders := make(map[net.Conn]*gob.Encoder)
	for {
		outgoing_msg := <-msg_channel
		clients_mu.Lock()
		for client := range clients {
			encoder, ok := encoders[client]
			if !ok {
				encoder = gob.NewEncoder(client)
				encoders[client] = encoder
			}
			err := encoder.Encode(&outgoing_msg)
			if err != nil {
				fmt.Println("Error Encoding: ", err.Error())
				delete(clients, client)
				client.Close()
			}
		}
		// Forget encoders for clients that have disconnected
		for client := range encoders {
			if _, ok := clients[client]; !ok {
				delete(encoders, client)
			}
		}
		clients_mu.Unlock()
	}
}
