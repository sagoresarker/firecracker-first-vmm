package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func InitMongoDB() {
	connectionString := "mongodb://user:pass@localhost:27021/firecrackerdb?authSource=admin&authMechanism=SCRAM-SHA-256"

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < 5; i++ {
		mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
		if err != nil {
			log.Println("Failed to create client:", err)
			time.Sleep(2 * time.Second)
			continue
		}
		err = mongoClient.Ping(ctx, nil)
		if err != nil {
			log.Println("Failed to ping:", err)
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	log.Println("Connected to MongoDB!")
}
