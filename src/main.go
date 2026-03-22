package main

import (
	"context"
	"log"
	"net/http"
	"os"
)

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

	http.HandleFunc("/api/lines", NewLinesHandler(store))

	// Serve static files from the src directory
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)

	log.Printf("Server starting on http://localhost:%s", port)
	log.Printf("Open http://localhost:%s in your browser", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
