package main

import (
	"cyberus/client-partner/internal/routes"
	"fmt"
	"net/http"
)

func main() {
	// Setup all routes
	routes.SetupRoutes()

	//Start the server on prod use port 8081
	fmt.Println("Starting cyberus-client-partner server on port 8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

}
