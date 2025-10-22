package storage

import (
	"bufio"
	"fmt"
	"kvstore/protocol"
	"os"
	"sync"
)

const LOGFILE = "data.aof"

type Store struct {
	data map[string]string
	mutex sync.RWMutex
	logFile *os.File
}

func NewStore() (*Store, error) {
	store := &Store{
		data: make(map[string]string),
	}

	if err := store.RecoverLog(); err != nil {
		return nil, fmt.Errorf("Error recovering log: %w\n", err)
	}

	logFile, err := os.OpenFile(LOGFILE, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error opening logfile: %w\n", err)
	}
	store.logFile = logFile

	return store, nil
}

func (s *Store) Close() (error) {
	if s.logFile != nil {
		return s.logFile.Close()
	}
	return nil
}

func (s *Store) Get(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	value, exists := s.data[key]
	return value, exists
}

func (s *Store) Set(key, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, err := s.logFile.WriteString(fmt.Sprintf("SET %s %s\n", key, value))
	if err != nil {
		return fmt.Errorf("failed to write to log: %w", err)
	}

	err = s.logFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to persist: %w", err)
	}

	s.data[key] = value
	return nil
}

func (s *Store) Delete(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, err := s.logFile.WriteString(fmt.Sprintf("DELETE %s\n", key))
	if err != nil {
		return fmt.Errorf("failed to write to log: %w", err)
	}

	err = s.logFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to persist: %w", err)
	}

	delete(s.data, key)
	return nil
}

func (s *Store) RecoverLog() error {
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
		
		op, key, value, err := protocol.ParseCommand(line)
		if err != nil {
			fmt.Printf("WARNING: skipping invalid line: %v\n", err)
			continue
		}

		switch op {
		case "SET":
			s.data[key] = value
		case "DELETE":
			delete(s.data, key)
		}
	}

	err = scanner.Err()
	if err != nil {
		return err
	}

	fmt.Printf("Recovery complete.\n")
	
	return nil
}