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

func GetAllProjectsByTenantID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	tenantID := params["tenantId"]

	collection := config.GetCollection("projects")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tenantIdBson, err := bson.ObjectIDFromHex(tenantID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "invalid tenant ID format",
			"msg":   err.Error(),
		})
		return
	}

	cursor, err := collection.Find(ctx, bson.M{"tenant_id": tenantIdBson})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to fetch projects",
			"msg":   err.Error(),
		})
		return
	}
	defer cursor.Close(ctx)

	var projects []models.Project
	if err := cursor.All(ctx, &projects); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to decode",
			"msg":   err.Error(),
		})
		return
	}
	if projects == nil {
		projects = []models.Project{}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(projects)
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	tenantID := params["tenantId"]

	var project models.Project

	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "invalid request payload",
			"msg":   err.Error(),
		})
		return
	}
	now := time.Now()
	project.CreatedAt = now
	project.UpdatedAt = now

	tenantObjectID, err := bson.ObjectIDFromHex(tenantID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "invalid tenant ID format",
			"msg":   err.Error(),
		})
		return
	}
	project.TenantID = tenantObjectID

	collection := config.GetCollection("projects")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, project)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to create project",
			"msg":   err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "project created successfully",
		"data":    result,
	})
}
