package handlers

import (
	"bytes"
	"context"
	"deployment-service/config"
	"deployment-service/models"
	"deployment-service/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// validateServiceExists checks if service exists in core-service
func validateServiceExists(serviceID string) (bool, error) {
	url := fmt.Sprintf("%s/services/%s", config.CoreServiceURL, serviceID)

	// Generate JWT token for authentication
	token, err := utils.GenerateServiceToken()
	if err != nil {
		return false, fmt.Errorf("failed to generate auth token: %v", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create request with JWT Bearer token
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to connect to core service at %s: %v", config.CoreServiceURL, err)
	}
	defer resp.Body.Close()

	// Service not found
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	// Authentication/Authorization errors
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("authentication failed with core service: %s", string(body))
	}

	// Unexpected error from core service
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("core service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Optionally validate response body to ensure it's a valid service
	var serviceResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&serviceResponse); err != nil {
		return false, fmt.Errorf("failed to decode core service response: %v", err)
	}

	// Verify the service has required fields
	if _, ok := serviceResponse["id"]; !ok {
		return false, fmt.Errorf("invalid service response: missing id field")
	}

	return true, nil
}

// callExternalAPI calls the external deployment API
func callExternalAPI(deployment models.Deployment) error {

	url := "https://jsonplaceholder.typicode.com/posts"

	payload := map[string]interface{}{
		"deployment_id": deployment.ID.Hex(),
		"service_id":    deployment.ServiceID.Hex(),
		"status":        deployment.Status,
		"timestamp":     time.Now().Unix(),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Create client with timeout
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to call external API: %v", err)
	}
	defer resp.Body.Close()

	// Accept both 200 and 201 status codes
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("external API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// updateDeploymentStatus updates the deployment status in the database
func updateDeploymentStatus(deploymentID bson.ObjectID, status models.DeploymentStatus) error {
	collection := config.GetCollection("deployments")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": deploymentID}, update)
	return err
}

// CreateDeployment creates a new deployment
func CreateDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	serviceID := params["serviceId"]

	// Validate serviceID format
	serviceObjectID, err := bson.ObjectIDFromHex(serviceID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "invalid service ID format",
			"msg":   err.Error(),
		})
		return
	}

	// Validate service exists in core-service
	exists, err := validateServiceExists(serviceID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to validate service",
			"msg":   err.Error(),
		})
		return
	}

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "service not found",
		})
		return
	}

	// Create deployment with Pending status
	now := time.Now()
	deployment := models.Deployment{
		ServiceID: serviceObjectID,
		Status:    models.StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	collection := config.GetCollection("deployments")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, deployment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to create deployment",
			"msg":   err.Error(),
		})
		return
	}

	// Get the inserted deployment ID
	insertedID := result.InsertedID.(bson.ObjectID)
	deployment.ID = insertedID

	// Process deployment asynchronously
	go func() {
		// Call external API
		if err := callExternalAPI(deployment); err != nil {
			// Mark as Failed
			updateDeploymentStatus(insertedID, models.StatusFailed)
			fmt.Printf("Deployment %s failed: %v\n", insertedID.Hex(), err)
			return
		}

		// Mark as Running
		if err := updateDeploymentStatus(insertedID, models.StatusRunning); err != nil {
			fmt.Printf("Failed to update deployment %s to Running: %v\n", insertedID.Hex(), err)
		} else {
			fmt.Printf("Deployment %s is now Running\n", insertedID.Hex())
		}
	}()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "deployment created successfully",
		"data": map[string]interface{}{
			"id":         insertedID,
			"service_id": serviceObjectID,
			"status":     models.StatusPending,
		},
	})
}

// GetDeploymentsByServiceID retrieves all deployments for a service
func GetDeploymentsByServiceID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	serviceID := params["serviceId"]

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

	serviceObjectID, err := bson.ObjectIDFromHex(serviceID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "invalid service ID format",
			"msg":   err.Error(),
		})
		return
	}

	collection := config.GetCollection("deployments")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"service_id": serviceObjectID}

	totalCount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to count deployments",
			"msg":   err.Error(),
		})
		return
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by newest first

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to fetch deployments",
			"msg":   err.Error(),
		})
		return
	}
	defer cursor.Close(ctx)

	var deployments []models.Deployment
	if err := cursor.All(ctx, &deployments); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "failed to decode deployments",
			"msg":   err.Error(),
		})
		return
	}

	if deployments == nil {
		deployments = []models.Deployment{}
	}

	totalPages := (int(totalCount) + limit - 1) / limit

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":        deployments,
		"page":        page,
		"limit":       limit,
		"total_count": totalCount,
		"total_pages": totalPages,
	})
}
