package handlers

import (
	"context"
	"core-service/config"
	"core-service/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func CreateService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	projectID := params["projectId"]

	var service models.Service

	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "invalid request payload",
			"msg":   err.Error(),
		})
		return
	}

	now := time.Now()
	service.CreatedAt = now
	service.UpdatedAt = now

	projectObjectID, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "invalid project ID format",
			"msg":   err.Error(),
		})
		return
	}

	service.ProjectID = projectObjectID

	collection := config.GetCollection("services")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, service)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to create service",
			"msg":   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "service created successfully",
		"data":    result,
	})
}

func GetAllServicesByProjectID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	projectID := params["projectId"]

	collection := config.GetCollection("services")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	projectIdBson, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "invalid project ID format",
			"msg":   err.Error(),
		})
		return
	}

	cursor, err := collection.Find(ctx, bson.M{"project_id": projectIdBson})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to fetch services",
			"msg":   err.Error(),
		})
		return
	}
	defer cursor.Close(ctx)

	var services []models.Service
	if err := cursor.All(ctx, &services); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to decode",
			"msg":   err.Error(),
		})
		return
	}

	if services == nil {
		services = []models.Service{}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services)
}

func GetServiceByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	serviceID := params["id"]

	serviceObjectID, err := bson.ObjectIDFromHex(serviceID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "invalid service ID format",
			"msg":   err.Error(),
		})
		return
	}

	collection := config.GetCollection("services")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var service models.Service
	err = collection.FindOne(ctx, bson.M{"_id": serviceObjectID}).Decode(&service)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "service not found",
			"msg":   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(service)
}
