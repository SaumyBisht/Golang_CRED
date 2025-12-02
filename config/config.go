package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client

const DB_URL string = "mongodb+srv://shantanubose_db_user:S92rdJWvGn50Dqoc@cluster0.mcfuxtt.mongodb.net/?appName=Cluster0"

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
	return MongoClient.Database("productdb").Collection(collectionName)
}
