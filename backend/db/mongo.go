package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func Connect() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI not set")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "workflow_approval"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}

	// Ping (เช็คว่า Atlas ต่อได้จริง)
	ctxPing, cancelPing := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelPing()

	if err := client.Ping(ctxPing, nil); err != nil {
		log.Fatal("Mongo ping error:", err)
	}

	Client = client
	Database = client.Database(dbName)

	log.Printf("✅ Mongo connected: db=%s", dbName)
}

func Disconnect() {
	if Client == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = Client.Disconnect(ctx)
}
