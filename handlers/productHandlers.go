// products.go

package handlers

import (
	"context"
	// "log"
	"modfile/db"
	"modfile/models"
	"net/http"
	// "os"

	"github.com/gin-gonic/gin"
	// "github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
)


// func InitCollection() {
// 	// Load the .env file
// 	if err := godotenv.Load(); err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	databaseName = os.Getenv("DATABASE_NAME")
// 	inventoryCollection = os.Getenv("INVENTORY_COLLECTION")

// 	// Initialize the collection from the global database client
// 	dbcollection = db.Client.Database(databaseName).Collection(inventoryCollection)
// }

func GetUserProducts(c *gin.Context) {
	// Get the user ID from the token
	userID, _ := c.Get("user")
	// fmt.Println("userId",userID)
	// Fetch products for the specific user
	cursor, err := db.Client.Database("Go-Gin-Inventory").Collection("inventory").Find(context.Background(), bson.M{"userID": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
		return
	}

	defer cursor.Close(context.Background())

	var products []models.InventoryItem
	if err := cursor.All(context.Background(), &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding products"})
		return
	}

	// fmt.Println("products",products)

	c.JSON(http.StatusOK, products)
}

func getProductById(c *gin.Context){
	id:=c.Param("id")
	// fmt.Println("id get product",id)
	//conver the id coming from api to the format that mongodb requires.
	objectId,err:=primitive.ObjectIDFromHex(id)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Invalid ID format",
		})
		return
	}

	//get userID from the token
	userID,_:=c.Get("user")
	// fmt.Println("userid",userID)

	var product models.InventoryItem
	err = dbcollection.FindOne(context.Background(),bson.M{"_id":objectId,"userID":userID}).Decode(&product)
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

	//get userID from token
	userID,_:=c.Get("user")
	// fmt.Println("userid in createproduct",userID)
	product.UserID=userID.(string)//convert ID to string

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
	// fmt.Println("id update product",id)
	objectId,err:=primitive.ObjectIDFromHex(id)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Invalid ID format",
		})
		return
	}

	//get userID from token
	userID,_:=c.Get("user")
	// fmt.Println("userid",userID)

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

	//Find the document with specific id and userID and update it
	result,err:=dbcollection.UpdateOne(
		context.Background(),
		bson.M{"_id":objectId,"userID":userID},
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

	//Get the userID from token
	userID,_:=c.Get("user")

	//delete the document with specified ID and userID
	result,err:=dbcollection.DeleteOne(context.Background(),bson.M{"_id":objectId,"userID":userID})
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
