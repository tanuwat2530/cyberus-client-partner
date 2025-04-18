package controllers

import (
	"CyberusGolangShareLibrary/utilities"
	services "cyberus/client-partner/internal/services"
	"net/http"
)

func ClientManagement(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := services.WapRedirectProcessRequest(r)

	utilities.ResponseWithJSON(w, http.StatusOK, response)
}
