package main

import (
    "facegram/database"
	"log"
	"github.com/gin-gonic/gin"
)

func main() {
    database.ConnectDB()
    
    r := gin.Default()

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

    log.Println("Server running at http://localhost:8080")
    if err := r.Run(); err != nil {
        log.Fatal("Failed to run server: ", err)
    }
}
