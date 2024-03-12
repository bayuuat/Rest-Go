package controller

import (
	"RestAPI/db"
	"RestAPI/helper"
	"RestAPI/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ItemsBody struct {
	ItemID      int64  `json:"item_id"`
	ItemCode    string `json:"item_code" binding:"required"`
	Description string `json:"description" binding:"required"`
	Quantity    int64  `json:"quantity" binding:"required"`
}

type OrderWithItemsBody struct {
	CustomerName string      `json:"customer_name" binding:"required"`
	OrderedAt    time.Time   `json:"ordered_at"`
	Items        []ItemsBody `json:"items" binding:"required"`
}

func GetOrder(ctx *gin.Context) {
	var orders []model.Order
	db.DB.Preload("Items").Find(&orders)

	ctx.JSON(http.StatusOK, orders)
}

func GetOrderById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, helper.ErrorResponse{Message: "invalid required param"})
		return
	}

	var orderWithItems model.Order
	db.DB.Preload("Items").First(&orderWithItems, "order_id = ?", id)

	result := map[string]interface{}{
		"customerName": orderWithItems.CustomerName,
		"orderedAt":    orderWithItems.OrderedAt,
		"items":        orderWithItems.Items,
	}

	ctx.JSON(http.StatusOK, result)
}

func CreateOrder(ctx *gin.Context) {
	var orderWithItems OrderWithItemsBody

	if err := ctx.BindJSON(&orderWithItems); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	order := model.Order{
		CustomerName: orderWithItems.CustomerName,
		OrderedAt:    orderWithItems.OrderedAt,
	}

	if err := db.DB.Create(&order).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	var items []model.Item
	for _, itemData := range orderWithItems.Items {
		items = append(items, model.Item{
			ItemCode:    itemData.ItemCode,
			Description: itemData.Description,
			Quantity:    itemData.Quantity,
			OrderID:     order.OrderID,
		})
	}

	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create items"})
		return
	}

	tx.Commit()

	ctx.JSON(http.StatusCreated, gin.H{"message": "Success create data order", "data": orderWithItems})
}

func UpdateOrder(ctx *gin.Context) {
	orderId, err := strconv.Atoi(ctx.Param("orderId"))

	var orderWithItems OrderWithItemsBody

	if orderId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, helper.ErrorResponse{Message: "invalid required param"})
		return
	}

	if err := ctx.ShouldBindJSON(&orderWithItems); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existOrder model.Order

	if err := db.DB.Table("orders").Where("order_id = ?", orderId).First(&existOrder).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order data not found"})
		return
	}

	existOrder.CustomerName = orderWithItems.CustomerName
	existOrder.OrderedAt = orderWithItems.OrderedAt

	if err := db.DB.Save(&existOrder).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	for _, itemData := range orderWithItems.Items {
		var existingItem model.Item

		if err := db.DB.Table("items").Where("item_id = ?", itemData.ItemID).First(&existingItem).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Item " + strconv.FormatInt(itemData.ItemID, 10) + " not found"})
			return
		}

		existingItem.ItemCode = itemData.ItemCode
		existingItem.Description = itemData.Description
		existingItem.Quantity = itemData.Quantity
		existingItem.OrderID = int64(orderId)

		if err := db.DB.Save(&existingItem).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
			return
		}
	}

	tx.Commit()

	ctx.JSON(http.StatusCreated, orderWithItems)
}

func DeleteOrder(ctx *gin.Context) {
	orderId, err := strconv.Atoi(ctx.Param("orderId"))

	if orderId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, helper.ErrorResponse{Message: "Invalid required param"})
		return
	}

	var order model.Order
	if err := db.DB.First(&order, orderId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
		return
	}

	if err := db.DB.Where("order_id = ?", orderId).Delete(&model.Item{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order items"})
		return
	}

	if err := db.DB.Delete(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order and associated items deleted successfully"})
}
