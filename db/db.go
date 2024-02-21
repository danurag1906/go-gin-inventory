package db

import (
	"context"
	"log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"github.com/joho/godotenv"
)

var Client *mongo.Client

func Init() {
	// Connect to MongoDB

	// Load the .env file
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

	databaseURL:=os.Getenv("MONGO_CONN_URL")

	clientOptions := options.Client().ApplyURI(databaseURL)
	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the MongoDB server to check the connection
	err = Client.Ping(context.Background(),nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")
}