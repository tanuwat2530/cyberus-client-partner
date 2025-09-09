package controllers

import (
	"CyberusGolangShareLibrary/utilities"
	"cyberus/client-partner/internal/services"
	"net/http"
)

func AdminDnDetailController(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := services.AdminDnDetailService(r)
	utilities.ResponseWithJSON(w, http.StatusOK, response)
}
