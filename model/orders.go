package model

import (
	"time"
)

type Order struct {
	OrderID      int64      `json:"order_id" gorm:"primaryKey;autoIncrement:true"`
	CustomerName string     `json:"customer_name"`
	OrderedAt    time.Time  `json:"ordered_at" gorm:"default:current_timestamp"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	Items        []Item
}
