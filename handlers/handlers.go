package handlers

import (
	"context"
	"modfile/db"
	"modfile/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var dbcollection *mongo.Collection

func InitCollection(){
	// Initialize the collection from the global database client
	dbcollection=db.Client.Database("Go-Gin-Inventory").Collection("inventory")
}

func SetupRoutes(r *gin.Engine){
	r.GET("/allProducts",getAllProducts)
	r.GET("/products/:id",getProductById)
	r.POST("/createProduct",createProduct)
	r.PUT("/updateProduct/:id",updateProduct)
	r.DELETE("/deleteProduct/:id",deleteProduct)
}

func getAllProducts(c *gin.Context){
	cursor,err:=dbcollection.Find(context.Background(),bson.M{})
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Error fetching products",
		})
		return
	}

	defer cursor.Close(context.Background())

	var products []models.InventoryItem
	if err:=cursor.All(context.Background(),&products); err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Error decoding products",
		})
		return
	}

	c.JSON(http.StatusOK,products)
}

func getProductById(c *gin.Context){
	id:=c.Param("id")
	objectId,err:=primitive.ObjectIDFromHex(id)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Invalid ID format",
		})
		return
	}

	var product models.InventoryItem
	err = dbcollection.FindOne(context.Background(),bson.M{"_id":objectId}).Decode(&product)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Error fetching product",
		})
		return
	}
	c.JSON(http.StatusOK,product)

}

func createProduct(c *gin.Context){
	var product models.InventoryItem
	if err:=c.ShouldBindJSON(&product); err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"invalid JSON format error",
		})
		return
	}
	result,err:=dbcollection.InsertOne(context.Background(),product)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"error creating product",
		})
		return
	}
	// Fetch the inserted product with details
	insertedProduct := models.InventoryItem{}
	err = dbcollection.FindOne(context.Background(), bson.M{"_id": result.InsertedID}).Decode(&insertedProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error fetching inserted product details",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"result": insertedProduct,
	})

}

func updateProduct(c *gin.Context){
	id:=c.Param("id")
	objectId,err:=primitive.ObjectIDFromHex(id)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Invalid ID format",
		})
		return
	}

	//binding the json to this variable
	var updateProduct models.InventoryItem
	if err:=c.ShouldBindJSON(&updateProduct); err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Invalid JSON Fomrat",
		})
		return
	}

	update:=bson.M{
		"$set":bson.M{
			"productName":updateProduct.ProductName,
			"units":updateProduct.Units,
			"price":updateProduct.Price,
		},
	}

	//Find the document with specific id and update it
	result,err:=dbcollection.UpdateOne(
		context.Background(),
		bson.M{"_id":objectId},
		update,
	)

	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Error updating product",
		})
		return
	}

	if result.ModifiedCount==0{
		c.JSON(http.StatusNotFound,gin.H{
			"error":"Prodcut not found",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"Prodcut updated successfully",
	})

}


func deleteProduct(c *gin.Context){
	id:=c.Param("id")
	objectId,err:=primitive.ObjectIDFromHex(id)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Invalid ID format",
		})
		return
	}

	//delete the document with specified ID
	result,err:=dbcollection.DeleteOne(context.Background(),bson.M{"_id":objectId})
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Error deleting product",
		})
		return
	}

	if result.DeletedCount==0{
		c.JSON(http.StatusNotFound,gin.H{
			"error":"Product not found",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"Product deleted successfully",
	})
}