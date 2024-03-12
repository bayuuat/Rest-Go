package request

import "time"

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
