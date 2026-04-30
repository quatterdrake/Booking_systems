package models

import (
	"time"

	"gorm.io/gorm"
)

type RoomType string

const (
	RoomTypeSingle  RoomType = "single"
	RoomTypeDouble  RoomType = "double"
	RoomTypeSuite   RoomType = "suite"
	RoomTypeDeluxe  RoomType = "deluxe"
)

type Room struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Number      string         `gorm:"uniqueIndex;not null" json:"number"`
	Type        RoomType       `gorm:"not null" json:"type"`
	Description string         `json:"description"`
	Price       float64        `gorm:"not null" json:"price"`
	Capacity    int            `gorm:"not null;default:1" json:"capacity"`
	IsAvailable bool           `gorm:"default:true" json:"is_available"`
	HotelID     uint           `gorm:"not null" json:"hotel_id"`
	Hotel       Hotel          `gorm:"foreignKey:HotelID" json:"hotel,omitempty"`
	Amenities   []Amenity      `gorm:"many2many:room_amenities" json:"amenities,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Request/Response DTOs

type CreateRoomRequest struct {
	Number      string   `json:"number" binding:"required"`
	Type        RoomType `json:"type" binding:"required,oneof=single double suite deluxe"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	Capacity    int      `json:"capacity" binding:"required,min=1"`
	HotelID     uint     `json:"hotel_id" binding:"required"`
	AmenityIDs  []uint   `json:"amenity_ids"`
}

type UpdateRoomRequest struct {
	Number      string   `json:"number"`
	Type        RoomType `json:"type" binding:"omitempty,oneof=single double suite deluxe"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"omitempty,gt=0"`
	Capacity    int      `json:"capacity" binding:"omitempty,min=1"`
	IsAvailable *bool    `json:"is_available"`
	AmenityIDs  []uint   `json:"amenity_ids"`
}

type RoomFilter struct {
	HotelID     uint     `form:"hotel_id"`
	Type        RoomType `form:"type"`
	MinPrice    float64  `form:"min_price"`
	MaxPrice    float64  `form:"max_price"`
	Capacity    int      `form:"capacity"`
	IsAvailable *bool    `form:"is_available"`
	Page        int      `form:"page,default=1"`
	Limit       int      `form:"limit,default=10"`
}
