package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve static files from the src directory
	fs := http.FileServer(http.Dir("./src"))
	http.Handle("/", fs)

	log.Printf("Server starting on http://localhost:%s", port)
	log.Printf("Open http://localhost:%s in your browser", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
