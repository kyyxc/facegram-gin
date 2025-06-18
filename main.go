package main

import (
    "facegram/database"
	"github.com/gin-gonic/gin"
    "facegram/controllers"
)

func main() {
    database.ConnectDB()
    
    r := gin.Default()

    r.POST("/api/v1/auth/register", controllers.Register)
    r.POST("/api/v1/auth/login", controllers.Login)

    r.Run(":9090")
}
