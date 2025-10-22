package main

import (
	"fmt"
	"os"
	"kvstore/storage"
	"kvstore/server"
	"os/signal"
	"syscall"
)

var store *storage.Store

func main() {
	var err error
	store, err = storage.NewStore()
	if err != nil {
		fmt.Printf("Error initializing store: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	svr := server.NewServer("127.0.0.1:8080", store)

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	errChannel := make(chan error, 1)
	// run the server in a goroutine (seperate thread)
	go func() {
		if err := svr.Start(); err != nil {
			errChannel <- err
		}
	}()

	select {
	case <-sigChannel:
		fmt.Println("\nShutdown signal received...")
	case err := <-errChannel:
		fmt.Printf("Server error: %v\n", err)
	}

	if err := svr.Stop(); err != nil {
		fmt.Printf("Error stopping server: %v\n", err)
	}

	if err := store.Close(); err != nil {
		fmt.Printf("Error closing store: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Shutdown complete")
}