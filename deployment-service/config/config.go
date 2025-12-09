package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client
var CoreServiceURL string
var JWTSecret []byte
var ServiceName string = "deployment-service"

const DB_URL string = "mongodb+srv://shantanubose_db_user:S92rdJWvGn50Dqoc@cluster0.mcfuxtt.mongodb.net/?appName=Cluster0"

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	CoreServiceURL = os.Getenv("CORE_SERVICE_URL")
	if CoreServiceURL == "" {
		CoreServiceURL = "http://localhost:8081"
	}
	log.Printf("Core Service URL: %s", CoreServiceURL)

	secretEnv := os.Getenv("JWT_SECRET")
	if secretEnv != "" {
		JWTSecret = []byte(secretEnv)
		log.Println("JWT secret loaded from environment")
	} else {
		log.Fatal("JWT_SECRET environment variable is required for deployment service")
	}
}

func ConnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(DB_URL)

	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal("Failed to connect Mongo:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		log.Fatal("Failed to ping Mongo:", err)
	}

	MongoClient = client
	log.Println("Connected to MongoDB!")
}

func GetCollection(collectionName string) *mongo.Collection {
	if MongoClient == nil {
		log.Fatal("MongoClient not initialized")
	}
	if collectionName == "" {
		log.Fatal("Collection name cannot be empty")
	}
	return MongoClient.Database("deployment_db").Collection(collectionName)
}
