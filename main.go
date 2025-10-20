package main

import (
	"fmt"
	"os"
	"net"
	"bufio"
	"strings"
	"sync"
)

var (
	store = make(map[string]string)
	mutex sync.RWMutex
	logFile *os.File
)

func main() {
	listener, errListener := net.Listen("tcp", "127.0.0.1:8080")
	logFile, errOpenFile := os.OpenFile("data.aof", os.WRONLY | os.O_APPEND | os.O_CREATE, 0644)

	if errOpenFile != nil {
		fmt.Println("ERROR (OpenFile)")
	}
	defer logFile.Close()

	if errListener != nil {
		fmt.Println("ERROR (Listen)\n")
		os.Exit(1)
	}

	for {
		connection, errAccept := listener.Accept()

		if errAccept != nil {
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
	mutex.RLock()
	defer mutex.RUnlock()
	value, exists := store[key]

	if exists {
		connection.Write([]byte(value + "\n"))
	} else {
		message := fmt.Sprintf("ERROR: Key '%s' not found.\n", key)
		connection.Write([]byte(message))
	}
}

func handleSET(key string, value string, connection net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	_, err := logFile.WriteString(fmt.Sprintf("SET %s %s\n", key, value))
	if err != nil {
		connection.Write([]byte("ERROR: Failed to persist\n"))
		return
	}

	store[key] = value
	connection.Write([]byte("OK\n"))
}

func handleDELETE(key string, connection net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	_, err := logFile.WriteString(fmt.Sprintf("DELETE %s\n", key))
	if err != nil {
		connection.Write([]byte("ERROR: Failed to persist\n"))
		return
	}

	delete(store, key)
	connection.Write([]byte("OK\n"))
}