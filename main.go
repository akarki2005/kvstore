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
	LOGFILE = "data.aof"
)

func main() {
	errRecoverLog := recoverLog()
	if errRecoverLog != nil {
		fmt.Println("Error recovering logfile.")
		os.Exit(1)
	}

	var errOpenFile error
	logFile, errOpenFile = os.OpenFile(LOGFILE, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
	if errOpenFile != nil {
		fmt.Println("ERROR (OpenFile)")
		os.Exit(1)
	}
	defer logFile.Close()

	listener, errListener := net.Listen("tcp", "127.0.0.1:8080")
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

func parseCommand(command string) (op, key, value string, err error) {
	words := strings.Fields(strings.TrimSpace(command))
	if len(words) == 0 {
		return "", "", "", fmt.Errorf("Empty command.")
	}

	op = strings.ToUpper(words[0])

	switch op {
	case "GET", "DELETE":
		if len(words) != 2 {
			return "", "", "", fmt.Errorf("Usage: %s [key]", op)
		}
		key = words[1]
	case "SET":
		if len(words) != 3 {
			return "", "", "", fmt.Errorf("Usage: SET [key] [value]")
		}
		key = words[1]
		value = words[2]
	default:
		return "", "", "", fmt.Errorf("Unknown operation: %s", op)
	}

	return op, key, value, nil
}

func handleConnection(connection net.Conn) {
	defer connection.Close()
	reader := bufio.NewReader(connection)

	for {
		message, readErr := reader.ReadString('\n')

		if readErr != nil {
			connection.Write([]byte("Error reading message."))
			return
		}

		op, key, value, err := parseCommand(message)
		if err != nil {
			connection.Write([]byte(err.Error() + "\n"))
			continue
		}

		switch op {
		case "GET":
			handleGET(key, connection)
		case "SET":
			handleSET(key, value, connection)
		case "DELETE":
			handleDELETE(key, connection)
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
		connection.Write([]byte("ERROR: Failed to write to log\n"))
		return
	}

	err = logFile.Sync()
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
		connection.Write([]byte("ERROR: Failed to write to log\n"))
		return
	}

	err = logFile.Sync()
	if err != nil {
		connection.Write([]byte("ERROR: Failed to persist\n"))
		return
	}

	delete(store, key)
	connection.Write([]byte("OK\n"))
}

func recoverLog() error {
	file, err := os.Open(LOGFILE)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		
		op, key, value, err := parseCommand(line)
		if err != nil {
			fmt.Printf("WARNING: skipping invalid line: %v\n", err)
			continue
		}

		switch op {
		case "SET":
			store[key] = value
		case "DELETE":
			delete(store, key)
		}
	}

	err = scanner.Err()
	if err != nil {
		return err
	}

	fmt.Printf("Recovery complete.\n")
	
	return nil
}