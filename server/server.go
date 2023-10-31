package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
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

func server() {
	msg_channel := make(chan message)
	clients := make(map[net.Conn]struct{})

	fmt.Println("Server Running...")
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error Listening: ", err.Error())
		os.Exit(1)
	}
	defer server.Close()

	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
	fmt.Println("Waiting on client...")

	go broadcast_clients(msg_channel, clients)

	for {
		// Wait for connections
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error Accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Client Connected\n")
		// Add client connection to map
		clients[connection] = struct{}{}
		// Start routine to handle incoming client messages
		go read_client(connection, msg_channel)
	}
}

// Handle incoming messages from clients
func read_client(connection net.Conn, msg_channel chan message) {
	defer connection.Close()
	for {
		var incoming_msg message
		decoder := gob.NewDecoder(connection)
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
func broadcast_clients(msg_channel chan message, clients map[net.Conn]struct{}) {
	for {
		outgoing_msg := <-msg_channel
		for client := range clients {
			encoder := gob.NewEncoder(client)
			err := encoder.Encode(&outgoing_msg)
			if err != nil {
				fmt.Println("Error Encoding: ", err.Error())
				delete(clients, client)
				break
			}
		}
	}
}
