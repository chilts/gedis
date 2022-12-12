package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Store is a simple in-memory storage for our key-value data
type Store struct {
	data map[string]string
}

// Set sets the value for a given key in the store
func (s *Store) Set(key, value string) error {
	s.data[key] = value
	return nil
}

// Get retrieves the value for a given key from the store
func (s *Store) Get(key string) (string, error) {
	value, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return value, nil
}

// Server is the main Redis server object
type Server struct {
	store *Store
}

// processRequest processes a given command and arguments, and returns the response
func (s *Server) processRequest(cmd string, args []string) (string, error) {
	switch cmd {
	case "set":
		if len(args) != 2 {
			return "", fmt.Errorf("invalid number of arguments for 'set' command")
		}
		return "OK", s.store.Set(args[0], args[1])
	case "get":
		if len(args) != 1 {
			return "", fmt.Errorf("invalid number of arguments for 'get' command")
		}
		value, err := s.store.Get(args[0])
		if err != nil {
			return "", err
		}
		return value, nil
	default:
		return "", fmt.Errorf("invalid command: %s", cmd)
	}
}

func main() {
	store := &Store{
		data: make(map[string]string),
	}
	server := &Server{store: store}

	// Hello world, the web server
	handler := func(w http.ResponseWriter, req *http.Request) {
		var input []string
		err := json.NewDecoder(req.Body).Decode(&input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Process this request
		// args := input[1:]

		// Process each request
		response, err := server.processRequest(input[0], input[1:])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Return the responses as JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.HandleFunc("/", handler)

	log.Println("Listing for requests at http://localhost:8000/")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
