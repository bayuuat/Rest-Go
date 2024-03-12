package routes

import (
	"RestAPI/controller"

	"github.com/gin-gonic/gin"
)

func OrderRoute(r *gin.RouterGroup) {
	r.GET("/orders", controller.GetOrder)
	r.POST("/orders", controller.CreateOrder)
	r.PUT("/orders/:orderId", controller.UpdateOrder)
}
