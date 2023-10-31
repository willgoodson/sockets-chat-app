package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"github.com/google/uuid"
	"net"
	"os"
	"strings"
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
	client()
}

func client() {
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	// Create a unique client id
	client_id := uuid.New().String()
	// Input Username
	fmt.Printf("Username: ")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	username = strings.Trim(username, "\n")
	// Start checking for messages from server
	go read_server(connection, client_id, username)

	for {
		// Prompt input
		fmt.Printf("%s> ", username)
		// Wait for input
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		outgoing_msg := message{Id: client_id, Username: username, Message: input}

		encoder := gob.NewEncoder(connection)
		err := encoder.Encode(&outgoing_msg)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
	}
}

func read_server(connection net.Conn, client_id string, username string) {
	defer connection.Close()
	var incoming_msg message
	for {
		decoder := gob.NewDecoder(connection)
		err := decoder.Decode(&incoming_msg)
		if err != nil {
			fmt.Println("Error Decoding: ", err.Error())
			break
		}
		if incoming_msg.Id != client_id {
			fmt.Printf("\n%s> %s%s>", incoming_msg.Username, incoming_msg.Message, username)
		}
	}
}
