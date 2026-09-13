package main

import (
	"context"
	"log"
	"net/http"
	"os"
)

func NewServerMux(store LineStore) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(sessionPagePathPrefix, NewSessionPageHandler())
	mux.HandleFunc(apiSessionsPathPrefix, NewSessionAPIHandler(store))
	mux.HandleFunc("/", NewRootHandler())

	return mux
}

func main() {
	ctx := context.Background()

	// Get port from environment variable or use default
	port := os.Getenv("SC_PORT")
	if port == "" {
		port = "8080"
	}

	store, err := NewLineStoreFromEnv(ctx)
	if err != nil {
		log.Fatal("Server failed to initialize store:", err)
	}

	server := NewServerMux(store)

	log.Printf("Server starting on http://localhost:%s", port)
	log.Printf("Open http://localhost:%s in your browser", port)

	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
