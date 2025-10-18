package main

import (
	"fmt"
	"os"
	"net"
	"bufio"
	"strings"
)

var store = make(map[string]string)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")

	if err != nil {
		fmt.Println("ERROR (Listen)\n")
		os.Exit(1)
	}

	for {
		connection, err := listener.Accept()

		if err != nil {
			fmt.Println("ERROR (Accept)\n")
			os.Exit(1)
		}

		go handleConnection(connection)
	}

}

func handleConnection(connection net.Conn) {
	defer connection.Close()
	reader := bufio.NewReader(connection)

	for {
		message, err := reader.ReadString('\n')

		if err != nil {
			connection.Write([]byte("Error reading message."))
			return
		}

		message = strings.TrimSpace(message)
		words := strings.Split(message, " ")

		if len(words) == 0 {
			connection.Write([]byte(""))
			fmt.Println("ERROR: No command specified.\n")
			continue
		}

		switch words[0] {
		case "GET":
			if len(words) != 2 {
				connection.Write([]byte("Usage: GET [key]\n"))
				continue
			}
			handleGET(words[1], connection)
		case "SET":
			if len(words) != 3 {
				connection.Write([]byte("Usage: SET [key] [value]\n"))
				continue
			}
			handleSET(words[1], words[2], connection)
		case "DELETE":
			if len(words) != 2 {
				connection.Write([]byte("Usage: DELETE [key]\n"))
				continue
			}
			handleDELETE(words[1], connection)
		default:
			connection.Write([]byte("Unknown operation.\n"))
		}
	}

}

func handleGET(key string, connection net.Conn) {
	value, exists := store[key]

	if exists {
		connection.Write([]byte(value + "\n"))
	} else {
		message := fmt.Sprintf("ERROR: Key '%s' not found.\n", key)
		connection.Write([]byte(message))
	}
}

func handleSET(key string, value string, connection net.Conn) {
	store[key] = value
	connection.Write([]byte("OK\n"))
}

func handleDELETE(key string, connection net.Conn) {
	delete(store, key)
	connection.Write([]byte("OK\n"))
}