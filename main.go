package main

import (
	"fmt"
	"os"
	"kvstore/storage"
	"kvstore/server"
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

	server := server.NewServer("127.0.0.1:8080", store)
	
	if err := server.Start(); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}

}