package model

import (
	"time"
)

type Item struct {
	ItemID      int64  `json:"item_id" gorm:"primaryKey;autoIncrement:true"`
	ItemCode    string `json:"item_code"`
	Description string `json:"description"`
	Quantity    int64  `json:"quantity"`
	OrderID     int64  `json:"order_id" gorm:"foreignKey:OrderID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
