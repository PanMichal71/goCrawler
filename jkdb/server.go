package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func storeHandler(db IDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		value, err := base64.StdEncoding.DecodeString(req.Value)
		if err != nil {
			http.Error(w, "Invalid base64 data", http.StatusBadRequest)
			return
		}
		fmt.Println("Server::storeHandler: Storing key: ", req.Key, " value: ", value)
		if err := db.Store(req.Key, value); err != nil {
			http.Error(w, "Error storing value", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func readHandler(db IDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		value, err := db.Read(key)
		if err != nil {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}
		response := struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{
			Key:   key,
			Value: base64.StdEncoding.EncodeToString(value),
		}
		json.NewEncoder(w).Encode(response)
	}
}

func deleteHandler(db IDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if err := db.Delete(key); err != nil {
			http.Error(w, "Error deleting key", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func existsHandler(db IDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		fmt.Println("Server::existsHandler: Checking existence of key: ", key)
		exists, err := db.Exists(key)
		if err != nil {
			http.Error(w, "Error checking existence", http.StatusInternalServerError)
			return
		}
		response := struct {
			Exists bool `json:"exists"`
		}{
			Exists: exists,
		}
		json.NewEncoder(w).Encode(response)
	}
}

func listKeysHandler(db IDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keys, err := db.ListKeys()
		if err != nil {
			http.Error(w, "Error listing keys", http.StatusInternalServerError)
			return
		}
		response := struct {
			Keys []string `json:"keys"`
		}{
			Keys: keys,
		}
		json.NewEncoder(w).Encode(response)
	}
}

func countHandler(db IDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count, err := db.Count()
		if err != nil {
			http.Error(w, "Error counting keys", http.StatusInternalServerError)
			return
		}
		response := struct {
			Count int `json:"count"`
		}{
			Count: count,
		}
		json.NewEncoder(w).Encode(response)
	}
}

func RunServer() {
	db := NewInMemoryDatabase()
	http.HandleFunc("/db/store", storeHandler(db))
	http.HandleFunc("/db/read", readHandler(db))
	http.HandleFunc("/db/delete", deleteHandler(db))
	http.HandleFunc("/db/exists", existsHandler(db))
	http.HandleFunc("/db/keys", listKeysHandler(db))
	http.HandleFunc("/db/count", countHandler(db))
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	http.ListenAndServe(":8080", nil)
}
