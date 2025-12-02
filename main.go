// main.go
package main

import (
	"log"
	"net/http"
	"products/config"
	"products/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Database connection
	config.ConnectDB()

	// Create router
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/products", handlers.GetAllProducts).Methods("GET")
	r.HandleFunc("/products/{id}", handlers.GetProductByID).Methods("GET")
	r.HandleFunc("/products", handlers.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", handlers.DeleteProduct).Methods("DELETE")

	// Start server
	log.Println("Server starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
