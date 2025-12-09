package main

import (
	"deployment-service/config"
	"deployment-service/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDB()

	r := mux.NewRouter()
	r.HandleFunc("/services/{serviceId}/deployments", handlers.CreateDeployment).Methods("POST")
	r.HandleFunc("/services/{serviceId}/deployments", handlers.GetDeploymentsByServiceID).Methods("GET")

	// Start server
	log.Println("Deployment Service starting on :8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}
