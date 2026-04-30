package models

import (
	"time"

	"gorm.io/gorm"
)

type ReservationStatus string

const (
	StatusPending   ReservationStatus = "pending"
	StatusConfirmed ReservationStatus = "confirmed"
	StatusCancelled ReservationStatus = "cancelled"
	StatusCompleted ReservationStatus = "completed"
)

type Reservation struct {
	ID          uint              `gorm:"primaryKey" json:"id"`
	UserID      uint              `gorm:"not null" json:"user_id"`
	User        User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	RoomID      uint              `gorm:"not null" json:"room_id"`
	Room        Room              `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	CheckIn     time.Time         `gorm:"not null" json:"check_in"`
	CheckOut    time.Time         `gorm:"not null" json:"check_out"`
	Guests      int               `gorm:"not null;default:1" json:"guests"`
	TotalPrice  float64           `gorm:"not null" json:"total_price"`
	Status      ReservationStatus `gorm:"default:'pending'" json:"status"`
	Notes       string            `json:"notes"`
	Payment     *Payment          `gorm:"foreignKey:ReservationID" json:"payment,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index" json:"-"`
}

type CreateReservationRequest struct {
	RoomID   uint      `json:"room_id" binding:"required"`
	CheckIn  time.Time `json:"check_in" binding:"required"`
	CheckOut time.Time `json:"check_out" binding:"required"`
	Guests   int       `json:"guests" binding:"required,min=1"`
	Notes    string    `json:"notes"`
}

type UpdateReservationRequest struct {
	CheckIn  *time.Time `json:"check_in"`
	CheckOut *time.Time `json:"check_out"`
	Guests   int        `json:"guests" binding:"omitempty,min=1"`
	Notes    string     `json:"notes"`
	Status   ReservationStatus `json:"status" binding:"omitempty,oneof=pending confirmed cancelled completed"`
}

type ReservationFilter struct {
	UserID   uint              `form:"user_id"`
	RoomID   uint              `form:"room_id"`
	Status   ReservationStatus `form:"status"`
	Page     int               `form:"page,default=1"`
	Limit    int               `form:"limit,default=10"`
}
