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

func GetAllTenants(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	collection := config.GetCollection("tenants")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to count tenants",
			"msg":   err.Error(),
		})
		return
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to fetch tenants",
			"msg":   err.Error(),
		})
		return
	}
	defer cursor.Close(ctx)

	var tenants []models.Tenant
	if err := cursor.All(ctx, &tenants); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to decode",
			"msg":   err.Error(),
		})
		return
	}
	if tenants == nil {
		tenants = []models.Tenant{}
	}

	totalPages := (int(totalCount) + limit - 1) / limit

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":        tenants,
		"page":        page,
		"limit":       limit,
		"total_count": totalCount,
		"total_pages": totalPages,
	})
}

func CreateTenant(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tenant models.Tenant

	if err := json.NewDecoder(r.Body).Decode(&tenant); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "invalid request payload",
			"msg":   err.Error(),
		})
		return
	}
	now := time.Now()
	tenant.CreatedAt = now
	tenant.UpdatedAt = now

	collection := config.GetCollection("tenants")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, tenant)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to create tenant",
			"msg":   err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "tenant created successfully",
		"data":    result,
	})
}

func GetTenantByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "invalid tenant ID format",
		})
		return
	}
	collection := config.GetCollection("tenants")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var tenant models.Tenant
	if err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&tenant); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to fetch tenant",
			"msg":   err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenant)

}
