package main

import (
	"log"
)

// main is the entry point of the application.
// It sets up the router and starts the HTTP server.
func main() {
	r := setupRouter()
	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
