package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	Category    string             `json:"category" bson:"category"`
	Stock       int                `json:"stock" bson:"stock"`
	CreatedBy   uint32             `json:"created_by" bson:"created_by"`
	UpdatedBy   uint32             `json:"updated_by" bson:"updated_by"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type CreateItemRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Category    string  `json:"category"`
	Stock       int     `json:"stock" binding:"min=0"`
}

type UpdateItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"min=0"`
	Category    string  `json:"category"`
	Stock       int     `json:"stock" binding:"min=0"`
}

type ItemResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    *Item  `json:"data,omitempty"`
}

type ItemsResponse struct {
	Message string  `json:"message"`
	Success bool    `json:"success"`
	Data    []*Item `json:"data"`
}
