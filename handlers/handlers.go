// package handlers

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"modfile/db"
// 	"modfile/models"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"golang.org/x/crypto/bcrypt"
// )

// var dbcollection *mongo.Collection


// func InitCollection(){
// 	// Load the .env file
//     if err := godotenv.Load(); err != nil {
//         log.Fatal("Error loading .env file")
//     }
// 	databaseName:=os.Getenv("DATABASE_NAME")
// 	inventoryCollection:=os.Getenv("INVENTORY_COLLECTION")
// 	// Initialize the collection from the global database client
// 	dbcollection=db.Client.Database(databaseName).Collection(inventoryCollection)
// }

// func SetupRoutes(r *gin.Engine){
// 	r.POST("/signup",signUp)
// 	r.POST("/signin",signIn)

// 	// Use the AuthMiddleware for routes that require authentication
// 	authGroup := r.Group("/auth")
// 	authGroup.Use(AuthMiddleware())
// 	{
// 		authGroup.GET("/allProducts", GetUserProducts)
// 		authGroup.GET("/products/:id", getProductById)
// 		authGroup.POST("/createProduct", createProduct)
// 		authGroup.PUT("/updateProduct/:id", updateProduct)
// 		authGroup.DELETE("/deleteProduct/:id", deleteProduct)
// 	}

// }

// func signUp(c *gin.Context){
// 	var user models.User
// 	if err:=c.ShouldBindJSON(&user); err!=nil{
// 		c.JSON(http.StatusBadRequest,gin.H{
// 			"error":"Invalid JSON format",
// 		})
// 		return
// 	}

// 	//check if username already exists
// 	existingUser:=models.User{}
// 	err:=db.Client.Database("Go-Gin-Inventory").Collection("users").FindOne(context.Background(),bson.M{"username":user.Username}).Decode(&existingUser)
// 	if err!=nil{
// 		c.JSON(http.StatusConflict,gin.H{
// 			"error":"Username alreay exists",
// 		})
// 		return
// 	}

// 	//Hash the password before storing it
// 	hashedPassword,err:=bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost)
// 	if err!=nil{
// 		c.JSON(http.StatusInternalServerError,gin.H{
// 			"error":"Error hashing pasword",
// 		})
// 		return
// 	}

// 	user.Password=string(hashedPassword)

// 	//insert new user
// 	_, err = db.Client.Database("Go-Gin-Inventory").Collection("users").InsertOne(context.Background(), user)
// 	if err!=nil{
// 		c.JSON(http.StatusInternalServerError,gin.H{
// 			"error":"Error creating user",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusCreated,gin.H{
// 		"message":"User created sucessfully!",
// 	})

// }


// func signIn(c *gin.Context){
// 	var user models.User
// 	if err:=c.ShouldBindJSON(&user);err!=nil{
// 		c.JSON(http.StatusBadRequest,gin.H{
// 			"error":"Invalid JSON format",
// 		})
// 		return
// 	}

// 	//database and collection name from .env
// 	databaseName:=os.Getenv("DATABASE_NAME")
// 	userCollection:=os.Getenv("USERs_COLLECTION")

// 	//Find the user by username
// 	storedUser:=models.User{}
// 	err:=db.Client.Database(databaseName).Collection(userCollection).FindOne(context.Background(),bson.M{"username":user.Username}).Decode(&storedUser)
// 	//decode is used to map the mongodb document to the user struct we defined in the code.
// 	if err!=nil{
// 		c.JSON(http.StatusUnauthorized,gin.H{
// 			"error":"Invalid credentials",
// 		})
// 		return
// 	}

// 	//valiate the hashed password using bcrypt
// 	if !validatePassword(user.Password,storedUser.Password){
// 		c.JSON(http.StatusUnauthorized,gin.H{
// 			"error":"Invalid Credentials",
// 		})
// 		return
// 	}

// 	//create a JWT token
// 	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
// 		"sub":storedUser.ID,
// 		"exp":time.Now().Add(time.Hour*24).Unix(),
// 		"iat":time.Now().Unix(),
// 		"username":storedUser.Username,
// 	})

// 	//get jwt secret from .dotenv
// 	jwtSceret:=os.Getenv("JWT_SECRET")

// 	//signin the token with secret key
// 	tokenString,err:=token.SignedString([]byte(jwtSceret))
// 	if err!=nil{
// 		c.JSON(http.StatusInternalServerError,gin.H{
// 			"error":"Error creating token",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK,gin.H{
// 		"token":tokenString,
// 		"username":storedUser.Username,
// 	})

// }

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 	   // Get the token from the Authorization header
// 		authHeader := c.GetHeader("Authorization")
// 		// fmt.Println("authheader",authHeader)
// 		// Check if the Authorization header is missing
// 		if authHeader == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 			return
// 		}

// 		// Split the header to get the actual token value
// 		tokenParts := strings.Split(authHeader, " ")
// 		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
// 			return
// 		}

// 		tokenString := tokenParts[1]

 
// 		//jwt secret from .env
// 		jwtSceret:=os.Getenv("JWT_SECRET")

// 	   // Parse the token
// 	   token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		  // Check the signing method
// 		  //if the signing method is ok it will return the secret used to sign the token
// 		  //if everthing is okay, we return the parsed token.
// 		  if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			 return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 		  }
// 		  return []byte(jwtSceret), nil
// 	   })

// 	   if err != nil {
// 		  c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		  return
// 	   }
 
// 	   // Check if the token is valid
// 	   if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		  // Set the user ID in the Gin context
// 		  //claims are the details which were used to make the token
// 		  //the sub i.e, userID is saved in the context and can be used by the next function.
// 		  c.Set("user", claims["sub"])
// 		  c.Next()
// 	   } else {
// 		  c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 	   }
// 	}
//  }
 

// // validatePassword function for validating bcrypt hashed password
// func validatePassword(inputPassword, storedPassword string) bool {
// 	// Compare bcrypt hashed password with plaintext input password
// 	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
// 	return err == nil
//  }


// // handlers/user.go

// func getAllProducts(c *gin.Context) {
// 	cursor, err := db.Client.Database("Go-Gin-Inventory").Collection("inventory").Find(context.Background(), bson.M{})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
// 		return
// 	}

// 	defer cursor.Close(context.Background())

// 	var products []models.InventoryItem
// 	if err := cursor.All(context.Background(), &products); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding products"})
// 		return
// 	}

// 	// Use sort.Reverse to reverse the order of the products slice
// 	// sort.Slice(products, func(i, j int) bool {
// 	// 	return i < j
// 	// })

// 	c.JSON(http.StatusOK, products)
// }

// func GetUserProducts(c *gin.Context) {
// 	// Get the user ID from the token
// 	userID, _ := c.Get("user")
// 	// fmt.Println("userId",userID)
// 	// Fetch products for the specific user
// 	cursor, err := db.Client.Database("Go-Gin-Inventory").Collection("inventory").Find(context.Background(), bson.M{"userID": userID})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
// 		return
// 	}

// 	defer cursor.Close(context.Background())

// 	var products []models.InventoryItem
// 	if err := cursor.All(context.Background(), &products); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding products"})
// 		return
// 	}

// 	// fmt.Println("products",products)

// 	c.JSON(http.StatusOK, products)
// }


// func getProductById(c *gin.Context){
// 	id:=c.Param("id")
// 	// fmt.Println("id get product",id)
// 	//conver the id coming from api to the format that mongodb requires.
// 	objectId,err:=primitive.ObjectIDFromHex(id)
// 	if err!=nil{
// 		c.JSON(http.StatusBadRequest,gin.H{
// 			"error":"Invalid ID format",
// 		})
// 		return
// 	}

// 	//get userID from the token
// 	userID,_:=c.Get("user")
// 	// fmt.Println("userid",userID)

// 	var product models.InventoryItem
// 	err = dbcollection.FindOne(context.Background(),bson.M{"_id":objectId,"userID":userID}).Decode(&product)
// 	if err!=nil{
// 		c.JSON(http.StatusInternalServerError,gin.H{
// 			"error":"Error fetching product",
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK,product)

// }

// func createProduct(c *gin.Context){
// 	var product models.InventoryItem

// 	//get userID from token
// 	userID,_:=c.Get("user")
// 	// fmt.Println("userid in createproduct",userID)
// 	product.UserID=userID.(string)//convert ID to string

// 	if err:=c.ShouldBindJSON(&product); err!=nil{
// 		c.JSON(http.StatusBadRequest,gin.H{
// 			"error":"invalid JSON format error",
// 		})
// 		return
// 	}
// 	result,err:=dbcollection.InsertOne(context.Background(),product)
// 	if err!=nil{
// 		c.JSON(http.StatusInternalServerError,gin.H{
// 			"error":"error creating product",
// 		})
// 		return
// 	}
// 	// Fetch the inserted product with details
// 	insertedProduct := models.InventoryItem{}
// 	err = dbcollection.FindOne(context.Background(), bson.M{"_id": result.InsertedID}).Decode(&insertedProduct)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "error fetching inserted product details",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"result": insertedProduct,
// 	})

// }

// func updateProduct(c *gin.Context){
// 	id:=c.Param("id")
// 	// fmt.Println("id update product",id)
// 	objectId,err:=primitive.ObjectIDFromHex(id)
// 	if err!=nil{
// 		c.JSON(http.StatusBadRequest,gin.H{
// 			"error":"Invalid ID format",
// 		})
// 		return
// 	}

// 	//get userID from token
// 	userID,_:=c.Get("user")
// 	// fmt.Println("userid",userID)

// 	//binding the json to this variable
// 	var updateProduct models.InventoryItem
// 	if err:=c.ShouldBindJSON(&updateProduct); err!=nil{
// 		c.JSON(http.StatusBadRequest,gin.H{
// 			"error":"Invalid JSON Fomrat",
// 		})
// 		return
// 	}

// 	update:=bson.M{
// 		"$set":bson.M{
// 			"productName":updateProduct.ProductName,
// 			"units":updateProduct.Units,
// 			"price":updateProduct.Price,
// 		},
// 	}

// 	//Find the document with specific id and userID and update it
// 	result,err:=dbcollection.UpdateOne(
// 		context.Background(),
// 		bson.M{"_id":objectId,"userID":userID},
// 		update,
// 	)

// 	if err!=nil{
// 		c.JSON(http.StatusInternalServerError,gin.H{
// 			"error":"Error updating product",
// 		})
// 		return
// 	}

// 	if result.ModifiedCount==0{
// 		c.JSON(http.StatusNotFound,gin.H{
// 			"error":"Prodcut not found",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK,gin.H{
// 		"message":"Prodcut updated successfully",
// 	})

// }


// func deleteProduct(c *gin.Context){
// 	id:=c.Param("id")
// 	objectId,err:=primitive.ObjectIDFromHex(id)
// 	if err!=nil{
// 		c.JSON(http.StatusBadRequest,gin.H{
// 			"error":"Invalid ID format",
// 		})
// 		return
// 	}

// 	//Get the userID from token
// 	userID,_:=c.Get("user")

// 	//delete the document with specified ID and userID
// 	result,err:=dbcollection.DeleteOne(context.Background(),bson.M{"_id":objectId,"userID":userID})
// 	if err!=nil{
// 		c.JSON(http.StatusInternalServerError,gin.H{
// 			"error":"Error deleting product",
// 		})
// 		return
// 	}

// 	if result.DeletedCount==0{
// 		c.JSON(http.StatusNotFound,gin.H{
// 			"error":"Product not found",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK,gin.H{
// 		"message":"Product deleted successfully",
// 	})
// }
package handlers

func main(){
	
}