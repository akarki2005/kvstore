package server

import (
	"bufio"
	"fmt"
	"kvstore/protocol"
	"kvstore/storage"
	"net"
)

type Server struct {
	address string
	store   *storage.Store
}

func NewServer(address string, store *storage.Store) *Server {
	return &Server{
		address: address,
		store:   store,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	fmt.Printf("Server listening on %s\n", s.address)

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go s.handleConnection(connection)
	}
}

func (s *Server) handleConnection(connection net.Conn) {
	defer connection.Close()
	reader := bufio.NewReader(connection)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		op, key, value, err := protocol.ParseCommand(message)
		if err != nil {
			connection.Write([]byte(err.Error() + "\n"))
			continue
		}

		switch op {
		case "GET":
			s.handleGET(key, connection)
		case "SET":
			s.handleSET(key, value, connection)
		case "DELETE":
			s.handleDELETE(key, connection)
		}
	}
}

func (s *Server) handleGET(key string, connection net.Conn) {
	value, exists := s.store.Get(key)
	if exists {
		connection.Write([]byte(value + "\n"))
	} else {
		message := fmt.Sprintf("ERROR: Key '%s' not found.\n", key)
		connection.Write([]byte(message))
	}
}

func (s *Server) handleSET(key, value string, connection net.Conn) {
	err := s.store.Set(key, value)
	if err != nil {
		connection.Write([]byte(fmt.Sprintf("ERROR: %v\n", err)))
		return
	}
	connection.Write([]byte("OK\n"))
}

func (s *Server) handleDELETE(key string, connection net.Conn) {
	err := s.store.Delete(key)
	if err != nil {
		connection.Write([]byte(fmt.Sprintf("ERROR: %v\n", err)))
		return
	}
	connection.Write([]byte("OK\n"))
}