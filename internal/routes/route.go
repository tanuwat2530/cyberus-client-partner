package routes

import (
	"cyberus/client-partner/internal/controllers"
	"net/http"
)

// SetupRoutes registers all application routes
func SetupRoutes() {
	// Register routes using http.HandleFunc
	http.HandleFunc("/client-service-api/add", controllers.AddServiceController)
	http.HandleFunc("/client-service-api/update", controllers.UpdateServiceController)
	http.HandleFunc("/client-api/add", controllers.AddClientController)
	http.HandleFunc("/client-api/update", controllers.UpdateClientController)
	http.HandleFunc("/api/", HomeHandler)
}

// HomeHandler for root endpoint
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to CYBERUS-CLIENT-PARTNER API power by GoLang ^_^"))
}
