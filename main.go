package main

import (
	"github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo"
	"modfile/handlers"
	"modfile/db"
	"github.com/gin-contrib/cors"
)	

var client *mongo.Client

const (
	ConnectionString ="mongodb+srv://gogin-inventory:gogin-inventory@cluster0.jdnkx7w.mongodb.net/?retryWrites=true&w=majority"
	DbName="go-inventory"
	CollectionName="inventory"
)


func main(){
	db.Init()
	handlers.InitCollection()
	r := gin.Default()

	config:=cors.DefaultConfig()
	config.AllowOrigins=[]string{"*"}
	r.Use(cors.New(config))

	handlers.SetupRoutes(r)
	r.Run(":8080")
}