package handlers

import (
	"net/http"
	"strconv"
	"time"

	"hotel-booking/config"
	"hotel-booking/models"
	"hotel-booking/utils"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct{}

func NewPaymentHandler() *PaymentHandler { return &PaymentHandler{} }

// CreatePayment godoc
// POST /payments
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req models.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	var reservation models.Reservation
	if err := config.DB.First(&reservation, req.ReservationID).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "reservation not found")
		return
	}

	if reservation.UserID != userID {
		utils.Fail(c, http.StatusForbidden, "access denied")
		return
	}

	if reservation.Status == models.StatusCancelled {
		utils.Fail(c, http.StatusBadRequest, "cannot pay for a cancelled reservation")
		return
	}

	// Check if payment already exists
	var existingPayment models.Payment
	if err := config.DB.Where("reservation_id = ?", req.ReservationID).First(&existingPayment).Error; err == nil {
		utils.Fail(c, http.StatusConflict, "payment already exists for this reservation")
		return
	}

	now := time.Now()
	payment := models.Payment{
		ReservationID: req.ReservationID,
		Amount:        reservation.TotalPrice,
		Method:        req.Method,
		Status:        models.PaymentStatusCompleted,
		TransactionID: req.TransactionID,
		PaidAt:        &now,
	}

	if err := config.DB.Create(&payment).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Confirm reservation after payment
	config.DB.Model(&reservation).Update("status", models.StatusConfirmed)

	utils.OK(c, http.StatusCreated, "payment successful", payment)
}

// GetPayment godoc
// GET /payments/:id
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	role := c.MustGet("userRole").(string)
	id, _ := strconv.Atoi(c.Param("id"))

	var payment models.Payment
	if err := config.DB.First(&payment, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "payment not found")
		return
	}

	// Check ownership
	if role != "admin" {
		var reservation models.Reservation
		config.DB.First(&reservation, payment.ReservationID)
		if reservation.UserID != userID {
			utils.Fail(c, http.StatusForbidden, "access denied")
			return
		}
	}

	utils.OK(c, http.StatusOK, "", payment)
}
