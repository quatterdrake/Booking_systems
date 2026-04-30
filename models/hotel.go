package models

import (
	"time"

	"gorm.io/gorm"
)

type Hotel struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Address     string         `gorm:"not null" json:"address"`
	City        string         `gorm:"not null" json:"city"`
	Country     string         `gorm:"not null" json:"country"`
	Stars       int            `gorm:"default:3" json:"stars"`
	Description string         `json:"description"`
	Phone       string         `json:"phone"`
	Email       string         `json:"email"`
	Rooms       []Room         `gorm:"foreignKey:HotelID" json:"rooms,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateHotelRequest struct {
	Name        string `json:"name" binding:"required,min=2"`
	Address     string `json:"address" binding:"required"`
	City        string `json:"city" binding:"required"`
	Country     string `json:"country" binding:"required"`
	Stars       int    `json:"stars" binding:"omitempty,min=1,max=5"`
	Description string `json:"description"`
	Phone       string `json:"phone"`
	Email       string `json:"email" binding:"omitempty,email"`
}

type UpdateHotelRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2"`
	Address     string `json:"address"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Stars       int    `json:"stars" binding:"omitempty,min=1,max=5"`
	Description string `json:"description"`
	Phone       string `json:"phone"`
	Email       string `json:"email" binding:"omitempty,email"`
}

type HotelFilter struct {
	City    string `form:"city"`
	Country string `form:"country"`
	Stars   int    `form:"stars"`
	Page    int    `form:"page,default=1"`
	Limit   int    `form:"limit,default=10"`
}
