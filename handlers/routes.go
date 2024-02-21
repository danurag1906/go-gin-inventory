// routes.go

package handlers

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/signup", signUp)
	r.POST("/signin", signIn)

	// Use the AuthMiddleware for routes that require authentication
	authGroup := r.Group("/auth")
	authGroup.Use(AuthMiddleware())
	{
		authGroup.GET("/allProducts", GetUserProducts)
		authGroup.GET("/products/:id", getProductById)
		authGroup.POST("/createProduct", createProduct)
		authGroup.PUT("/updateProduct/:id", updateProduct)
		authGroup.DELETE("/deleteProduct/:id", deleteProduct)
	}
}
