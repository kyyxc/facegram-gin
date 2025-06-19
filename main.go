package main

import (
	"facegram/database"
	"facegram/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	r := gin.Default()

	routes.SetupRoute(r)

	r.Run(":9090")
}
