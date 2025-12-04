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

	// Routes
	r.HandleFunc("/products", handlers.GetAllProducts).Methods("GET")
	r.HandleFunc("/products/{id}", handlers.GetProductByID).Methods("GET") // particular product
	r.HandleFunc("/products", handlers.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", handlers.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", handlers.DeleteProduct).Methods("DELETE")

	// tenant routes
	r.HandleFunc("/tenants", handlers.CreateTenant).Methods("POST")
	r.HandleFunc("/tenants", handlers.GetAllTenants).Methods("GET")
	r.HandleFunc("/tenants/{id}", handlers.GetTenantByID).Methods("GET")

	// Start server
	log.Println("Server starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
