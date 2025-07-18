package routes

import (
	"cyberus/client-partner/internal/controllers"
	"net/http"
)

// SetupRoutes registers all application routes
func SetupRoutes() {
	// Register routes using http.HandleFunc

	//ADD , DELETE , UPDATE , SERVICE FOR PARTNER
	http.HandleFunc("/client-service-api/list-service", controllers.ListServiceController)
	http.HandleFunc("/client-service-api/list-client", controllers.ListClientController)
	http.HandleFunc("/client-service-api/list-log", controllers.ListLogController)
	http.HandleFunc("/client-service-api/add", controllers.AddServiceController)
	http.HandleFunc("/client-service-api/update", controllers.UpdateServiceController)
	http.HandleFunc("/client-service-api/delete", controllers.DeleteServiceController)

	//ADD , DELETE , UPDATE PARTNER
	http.HandleFunc("/client-api/login", controllers.LoginClientController)
	http.HandleFunc("/client-api/session", controllers.SessionLoginController)
	http.HandleFunc("/client-api/add", controllers.AddClientController)
	http.HandleFunc("/client-api/update", controllers.UpdateClientController)

	//CLIENT REPORT
	http.HandleFunc("/client-api/login-report", controllers.LoginReportController)
	http.HandleFunc("/client-report/chart", controllers.ClientReportChartController)
	http.HandleFunc("/client-report/overall", controllers.ClientReportOverallController)
	http.HandleFunc("/client-report/shortcode-client", controllers.ListShortcodeClientController)
	http.HandleFunc("/client-report/request-history", controllers.RequestHistoryController)

	//ADMIN REPORT
	//http.HandleFunc("/admin-report/pie-report", controllers.AdminReportPieController)
	http.HandleFunc("/admin-report/all-shortcode", controllers.ListShortcodeAllController)

	http.HandleFunc("/api/", HomeHandler)
}

// HomeHandler for root endpoint
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to CYBERUS-CLIENT-PARTNER API power by GoLang ^_^"))
}
