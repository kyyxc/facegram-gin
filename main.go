package main

import (
	"facegram/controllers"
	"facegram/database"
	"facegram/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	r := gin.Default()

	r.POST("/api/v1/auth/register", controllers.Register)
	r.POST("/api/v1/auth/login", controllers.Login)

	auth := r.Group("/")
	auth.Use(middlewares.AuthMiddleware())
	auth.POST("/api/v1/posts", controllers.CreatePost)
	auth.DELETE("/api/v1/posts/:id", controllers.DeletePost)

	r.Run(":9090")
}
