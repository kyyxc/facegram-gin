package routes

import (
	"facegram/controllers"
	"facegram/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoute(route *gin.Engine) {
	api := route.Group("/api/v1")
	{
		// Auth
		api.POST("/auth/register", controllers.Register)
		api.POST("/auth/login", controllers.Login)

		// Protected route
		auth := api.Group("/")
		auth.Use(middlewares.AuthMiddleware())

		// Post
		auth.POST("/posts", controllers.CreatePost)
		auth.DELETE("/posts/:id", controllers.DeletePost)
		auth.GET("/posts", controllers.GetPost)

		// Follow
		auth.POST("/users/:username/follow", controllers.Follow)
	}
}
