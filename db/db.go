package db

import (
	"context"
	"log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Init() {
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb+srv://gogin-inventory:gogin-inventory@cluster0.jdnkx7w.mongodb.net/?retryWrites=true&w=majority")
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