package handlers

import (
	"context"
	"core-service/config"
	"core-service/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func GetAllProjectsByTenantID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	tenantID := params["tenantId"]

	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // Default limit
	}

	skip := (page - 1) * limit

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

	filter := bson.M{"tenant_id": tenantIdBson}

	totalCount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to count projects",
			"msg":   err.Error(),
		})
		return
	}

	// pagination parameters
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, findOptions)
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

	totalPages := (int(totalCount) + limit - 1) / limit

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":        projects,
		"page":        page,
		"limit":       limit,
		"total_count": totalCount,
		"total_pages": totalPages,
	})
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
