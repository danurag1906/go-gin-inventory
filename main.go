package main

import (
	"github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo"
	"modfile/handlers"
	"modfile/db"
	"github.com/gin-contrib/cors"
)	

var client *mongo.Client


func main(){
	db.Init()
	handlers.InitCollection()
	r := gin.Default()

	// Apply CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} // Add your React app's origin
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	r.Use(cors.New(config))

	handlers.SetupRoutes(r)
	r.Run(":8080")
}