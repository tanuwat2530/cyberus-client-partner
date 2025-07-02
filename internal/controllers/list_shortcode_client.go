package controllers

import (
	"CyberusGolangShareLibrary/utilities"
	services "cyberus/client-partner/internal/services"
	"net/http"
)

func ListShortcodeClientController(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := services.ListShortcodeClientService(r)
	utilities.ResponseWithJSON(w, http.StatusOK, response)
}
