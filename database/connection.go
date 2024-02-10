package database

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func DB() *mongo.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file...")
	}

	uri := os.Getenv("MONGO_URI")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	defer cancel()

	fmt.Println("Connected to mongodb...")

	return client
}

var Client *mongo.Client = DB()

func OpenCollection(client *mongo.Client, name string) *mongo.Collection {
	collection := client.Database("cluster0").Collection(name)
	return collection
}
