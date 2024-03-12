package routes

import (
	"RestAPI/controller"

	"github.com/gin-gonic/gin"
)

func ItemsRoute(r *gin.RouterGroup) {
	r.GET("/items", controller.GetItems)
	r.POST("/items", controller.CreateItem)
}
