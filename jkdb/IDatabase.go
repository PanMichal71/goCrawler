package main

import (
	"errors"
	"fmt"
	"sync"
)

type IDatabase interface {
	Store(key string, value []byte) error
	Read(key string) ([]byte, error)
	Delete(key string) error
	Exists(key string) (bool, error)
	ListKeys() ([]string, error)
	Count() (int, error)
}

type InMemoryDatabase struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		data: make(map[string][]byte),
	}
}

func (db *InMemoryDatabase) Store(key string, value []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	fmt.Println("Storing key: ", key, " value: ", value)
	db.data[key] = value
	return nil
}

func (db *InMemoryDatabase) Read(key string) ([]byte, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	fmt.Println("Reading key: ", key)
	value, exists := db.data[key]
	if !exists {
		return nil, errors.New("key not found")
	}
	return value, nil
}

func (db *InMemoryDatabase) Delete(key string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, exists := db.data[key]; !exists {
		return errors.New("key not found")
	}
	delete(db.data, key)
	return nil
}

func (db *InMemoryDatabase) Exists(key string) (bool, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	_, exists := db.data[key]

	fmt.Println("Checking existence of key: ", key, " exists: ", exists)
	return exists, nil
}

func (db *InMemoryDatabase) ListKeys() ([]string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	keys := make([]string, 0, len(db.data))
	for key := range db.data {
		keys = append(keys, key)
	}
	return keys, nil
}

func (db *InMemoryDatabase) Count() (int, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return len(db.data), nil
}
