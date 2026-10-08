package main

import (
	"go-postgres-api/config"
	"go-postgres-api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/items", handlers.GetAllItems)
		api.GET("/items/:id", handlers.GetItemByID)
		api.POST("/items", handlers.CreateItem)
		api.PUT("/items/:id", handlers.UpdateItem)
		api.DELETE("/items/:id", handlers.DeleteItem)
	}

	r.Run(":8081")
}
