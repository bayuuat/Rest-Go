package service

import (
	"RestAPI/db"
	"RestAPI/model"
	"RestAPI/request"
)

type OrderService struct{}

func (os *OrderService) CreateOrder(orderWithItems request.OrderWithItemsBody) (*model.Order, error) {
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
		return nil, err
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
		return nil, err
	}

	tx.Commit()

	return &order, nil
}

func (os *OrderService) UpdateOrder(orderID uint, orderWithItems request.OrderWithItemsBody) (*model.Order, error) {
	tx := db.DB.Begin()

	var existOrder model.Order

	if err := tx.Table("orders").Where("order_id = ?", orderID).First(&existOrder).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	existOrder.CustomerName = orderWithItems.CustomerName
	existOrder.OrderedAt = orderWithItems.OrderedAt

	if err := tx.Save(&existOrder).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, itemData := range orderWithItems.Items {
		var existingItem model.Item

		if err := tx.Table("items").Where("item_id = ?", itemData.ItemID).First(&existingItem).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		existingItem.ItemCode = itemData.ItemCode
		existingItem.Description = itemData.Description
		existingItem.Quantity = itemData.Quantity
		existingItem.OrderID = int64(orderID)

		if err := tx.Save(&existingItem).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	db.DB.Table("orders").Preload("Items").Where("order_id = ?", orderID).First(&existOrder)

	return &existOrder, nil
}

func (os *OrderService) DeleteOrder(orderID uint) error {
	var order model.Order
	if err := db.DB.First(&order, orderID).Error; err != nil {
		return err
	}

	if err := db.DB.Where("order_id = ?", orderID).Delete(&model.Item{}).Error; err != nil {
		return err
	}

	if err := db.DB.Delete(&order).Error; err != nil {
		return err
	}

	return nil
}
