package controller

import (
	"RestAPI/db"
	"RestAPI/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetItems(ctx *gin.Context) {
	var items []model.Item
	db.DB.Find(&items)
	ctx.JSON(http.StatusOK, items)
}

func CreateItem(c *gin.Context) {
	var newItem model.Item
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db.DB.Create(&newItem)
	c.JSON(http.StatusCreated, newItem)
}
