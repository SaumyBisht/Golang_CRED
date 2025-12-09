// main.go
package main

import (
	"core-service/config"
	"core-service/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Database connection
	config.ConnectDB()

	// Create router
	r := mux.NewRouter()

	// tenant routes
	r.HandleFunc("/tenants", handlers.CreateTenant).Methods("POST")
	r.HandleFunc("/tenants", handlers.GetAllTenants).Methods("GET")
	r.HandleFunc("/tenants/{id}", handlers.GetTenantByID).Methods("GET")

	// project routes
	r.HandleFunc("/tenants/{tenantId}/projects", handlers.CreateProject).Methods("POST")
	r.HandleFunc("/tenants/{tenantId}/projects", handlers.GetAllProjectsByTenantID).Methods("GET")

	//service routes
	r.HandleFunc("/projects/{projectId}/services", handlers.GetAllServicesByProjectID).Methods("GET")
	r.HandleFunc("/projects/{projectId}/services", handlers.CreateService).Methods("POST")

	// Start server
	log.Println("Server starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
