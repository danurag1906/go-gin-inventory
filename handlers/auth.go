// auth.go

package handlers

import (
	"context"
	"fmt"
	"log"
	"modfile/db"
	"modfile/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var dbcollection *mongo.Collection
var databaseName string
var inventoryCollection string
var usersCollection string
var jwtSecret string

func InitCollection() {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	databaseName = os.Getenv("DATABASE_NAME")
	inventoryCollection = os.Getenv("INVENTORY_COLLECTION")
	usersCollection=os.Getenv("USERS_COLLECTION")
	jwtSecret=os.Getenv("JWT_SECRET")
	// Initialize the collection from the global database client
	dbcollection = db.Client.Database(databaseName).Collection(inventoryCollection)
}

func signUp(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON format",
		})
		return
	}

	// Check if username already exists
	existingUser := models.User{}
	err := db.Client.Database(databaseName).Collection(usersCollection).FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&existingUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Username not found, proceed with signup
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error checking username existence",
			})
			return
		}
	} else {
		// Username already exists, return conflict response
		c.JSON(http.StatusConflict, gin.H{
			"error": "Username already exists",
		})
		return
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error hashing password",
		})
		return
	}

	user.Password = string(hashedPassword)

	// Insert new user
	_, err = db.Client.Database(databaseName).Collection(usersCollection).InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creating user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully!",
	})
}


func signIn(c *gin.Context){
	var user models.User
	if err:=c.ShouldBindJSON(&user);err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Invalid JSON format",
		})
		return
	}

	//Find the user by username
	storedUser:=models.User{}
	err:=db.Client.Database(databaseName).Collection(usersCollection).FindOne(context.Background(),bson.M{"username":user.Username}).Decode(&storedUser)
	//decode is used to map the mongodb document to the user struct we defined in the code.
	if err!=nil{
		c.JSON(http.StatusUnauthorized,gin.H{
			"error":"Invalid credentials",
		})
		return
	}

	//valiate the hashed password using bcrypt
	if !validatePassword(user.Password,storedUser.Password){
		c.JSON(http.StatusUnauthorized,gin.H{
			"error":"Invalid Credentials",
		})
		return
	}

	//create a JWT token
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"sub":storedUser.ID,
		"exp":time.Now().Add(time.Hour*24).Unix(),
		"iat":time.Now().Unix(),
		"username":storedUser.Username,
	})

	//signin the token with secret key
	tokenString,err:=token.SignedString([]byte(jwtSecret))
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Error creating token",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"token":tokenString,
		"username":storedUser.Username,
	})

}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
	   // Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		// fmt.Println("authheader",authHeader)
		// Check if the Authorization header is missing
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Split the header to get the actual token value
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		tokenString := tokenParts[1]

	   // Parse the token
	   token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		  // Check the signing method
		  //if the signing method is ok it will return the secret used to sign the token
		  //if everthing is okay, we return the parsed token.
		  if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			 return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		  }
		  return []byte(jwtSecret), nil
	   })

	   if err != nil {
		  c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		  return
	   }
 
	   // Check if the token is valid
	   if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		  // Set the user ID in the Gin context
		  //claims are the details which were used to make the token
		  //the sub i.e, userID is saved in the context and can be used by the next function.
		  c.Set("user", claims["sub"])
		  c.Next()
	   } else {
		  c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	   }
	}
 }


// validatePassword function for validating bcrypt hashed password
func validatePassword(inputPassword, storedPassword string) bool {
	// Compare bcrypt hashed password with plaintext input password
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
	return err == nil
}


