package models

import (
	"time"

	"gorm.io/gorm"
)

type PaymentStatus string
type PaymentMethod string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"

	PaymentMethodCard   PaymentMethod = "card"
	PaymentMethodCash   PaymentMethod = "cash"
	PaymentMethodOnline PaymentMethod = "online"
)

type Payment struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ReservationID uint           `gorm:"uniqueIndex;not null" json:"reservation_id"`
	Amount        float64        `gorm:"not null" json:"amount"`
	Method        PaymentMethod  `gorm:"not null" json:"method"`
	Status        PaymentStatus  `gorm:"default:'pending'" json:"status"`
	TransactionID string         `json:"transaction_id"`
	PaidAt        *time.Time     `json:"paid_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreatePaymentRequest struct {
	ReservationID uint          `json:"reservation_id" binding:"required"`
	Method        PaymentMethod `json:"method" binding:"required,oneof=card cash online"`
	TransactionID string        `json:"transaction_id"`
}
