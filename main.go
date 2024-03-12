package main

import (
	"RestAPI/db"
	"RestAPI/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	r := gin.Default()

	api := r.Group("/api")
	routes.ItemsRoute(api)
	routes.OrderRoute(api)

	r.Run(":8080")
}
