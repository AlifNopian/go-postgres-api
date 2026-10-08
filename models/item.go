package models

import "time"

type Item struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Price     float64   `json:"price" binding:"required,gt=0"`
	Stock     int       `json:"stock" binding:"required,gte=0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateItemInput struct {
	Name  *string  `json:"name"`
	Price *float64 `json:"price"`
	Stock *int     `json:"stock"`
}
