package models

import (
	"time"

	"gorm.io/gorm"
)

type Amenity struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"uniqueIndex;not null" json:"name"`
	Icon      string         `json:"icon"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateAmenityRequest struct {
	Name string `json:"name" binding:"required,min=2"`
	Icon string `json:"icon"`
}

// FavoriteRoom represents practicum 6 requirement
type FavoriteRoom struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	RoomID    uint      `gorm:"primaryKey" json:"room_id"`
	Room      Room      `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
