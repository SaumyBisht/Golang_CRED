package main

import (
	"core-service/config"
	"core-service/handlers"
	"core-service/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Database connection
	config.ConnectDB()

	// Create router
	r := mux.NewRouter()

	// Public routes (no authentication required)
	publicRouter := r.PathPrefix("").Subrouter()
	publicRouter.HandleFunc("/tenants", handlers.CreateTenant).Methods("POST")
	publicRouter.HandleFunc("/tenants", handlers.GetAllTenants).Methods("GET")
	publicRouter.HandleFunc("/tenants/{id}", handlers.GetTenantByID).Methods("GET")
	publicRouter.HandleFunc("/tenants/{tenantId}/projects", handlers.CreateProject).Methods("POST")
	publicRouter.HandleFunc("/tenants/{tenantId}/projects", handlers.GetAllProjectsByTenantID).Methods("GET")
	publicRouter.HandleFunc("/projects/{projectId}/services", handlers.GetAllServicesByProjectID).Methods("GET")
	publicRouter.HandleFunc("/projects/{projectId}/services", handlers.CreateService).Methods("POST")

	// Internal routes (require API key authentication)
	internalRouter := r.PathPrefix("").Subrouter()
	internalRouter.Use(middleware.AuthMiddleware)
	internalRouter.HandleFunc("/services/{id}", handlers.GetServiceByID).Methods("GET")

	// Start server
	log.Println("Server starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
