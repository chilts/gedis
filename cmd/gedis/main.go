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

type setResponse struct {
	Ok bool `json:"ok"`
}

func main() {
	store := &Store{
		data: make(map[string]string),
	}
	server := &Server{store: store}

	handleSet := func(w http.ResponseWriter, r *http.Request) {
		// Check that the request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Parse the request body as JSON
		var body interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Set the content type of the response to application/json
		w.Header().Set("Content-Type", "application/json")

		// Write the JSON response indicating that the operation was successful
		json.NewEncoder(w).Encode(setResponse{Ok: true})
	}

	// Hello world, the web server
	handler := func(w http.ResponseWriter, req *http.Request) {
		var input []string
		err := json.NewDecoder(req.Body).Decode(&input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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

	http.HandleFunc("/set", handleSet)
	http.HandleFunc("/", handler)

	log.Println("Listing for requests at http://localhost:8000/")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
