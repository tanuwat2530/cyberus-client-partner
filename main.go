package main

import (
	"cyberus/client-partner/internal/routes"
	"fmt"
	"net/http"
)

func main() {
	// Setup all routes
	routes.SetupRoutes()

	// Start the server on port 8080
	fmt.Println("Starting cyberus-client-partner server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

}
