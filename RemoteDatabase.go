package main

// an implementation of IDatabase interface, which stores data in a remote database
// and uses REST API to interact with it.
// APIs are:
// POST /db/store
// GET /db/read
// DELETE /db/delete
// GET /db/exists
// GET /db/keys
// GET /db/count

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type RemoteDatabase struct {
	url string
	mu  sync.RWMutex
}

func NewRemoteDatabase(url string) *RemoteDatabase {
	return &RemoteDatabase{
		url: url,
	}
}

func (db *RemoteDatabase) Store(key string, value []byte) error {

	db.mu.Lock()
	defer db.mu.Unlock()

	fmt.Println("Storing key: ", key, " value: ", value)

	req := struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{
		Key:   key,
		Value: base64.StdEncoding.EncodeToString(value), // Base64 encode the value
	}
	reqData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(db.url+"/db/store", "application/json", bytes.NewBuffer(reqData))
	if err != nil {
		fmt.Println("Error storing value: ", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error storing value: ", resp.Status)
		return errors.New("error storing value")
	}
	return nil
}

func (db *RemoteDatabase) Read(key string) ([]byte, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	resp, err := http.Get(db.url + "/db/read?key=" + key)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("key not found")
	}

	var respData struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	// Decode the response body to get the key-value pair
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	// Base64 decode the value
	decodedValue, err := base64.StdEncoding.DecodeString(respData.Value)
	if err != nil {
		return nil, err
	}

	return decodedValue, nil
}

func (db *RemoteDatabase) Delete(key string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	req, err := http.NewRequest(http.MethodDelete, db.url+"/db/delete?key="+key, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("error deleting key")
	}
	return nil
}

func (db *RemoteDatabase) Exists(key string) (bool, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	fmt.Println("Checking existence of key: ", key)
	resp, err := http.Get(db.url + "/db/exists?key=" + key)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	// print out resp
	var respData struct {
		Exists bool `json:"exists"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return false, err
	}

	fmt.Println("Key exists: ", respData.Exists)

	return respData.Exists, nil
}

func (db *RemoteDatabase) ListKeys() ([]string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	resp, err := http.Get(db.url + "/db/keys")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var keys []string
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return nil, err
	}
	return keys, nil
}

func (db *RemoteDatabase) Count() (int, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	resp, err := http.Get(db.url + "/db/count")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var respData struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return 0, err
	}
	return respData.Count, nil
}
