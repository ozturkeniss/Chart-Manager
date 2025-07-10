package model

import (
	"time"

	"gorm.io/gorm"
)

type Inventory struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ItemID    string         `json:"item_id" gorm:"uniqueIndex;not null"`
	Stock     int            `json:"stock" gorm:"not null;default:0"`
	MinStock  int            `json:"min_stock" gorm:"default:0"`
	MaxStock  int            `json:"max_stock" gorm:"default:1000"`
	UpdatedBy uint32         `json:"updated_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UpdateStockRequest struct {
	NewStock int `json:"new_stock" binding:"required,min=0"`
}

type StockResponse struct {
	Message string     `json:"message"`
	Success bool       `json:"success"`
	Data    *Inventory `json:"data,omitempty"`
}

type ItemsResponse struct {
	Message string       `json:"message"`
	Success bool         `json:"success"`
	Data    []*Inventory `json:"data"`
}
