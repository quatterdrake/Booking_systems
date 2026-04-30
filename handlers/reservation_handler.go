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

type ReservationHandler struct{}

func NewReservationHandler() *ReservationHandler { return &ReservationHandler{} }

// GetReservations godoc
// GET /reservations
func (h *ReservationHandler) GetReservations(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	role := c.MustGet("userRole").(string)

	var filter models.ReservationFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	query := config.DB.Model(&models.Reservation{})

	// Non-admins only see their own reservations
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	} else if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.RoomID > 0 {
		query = query.Where("room_id = ?", filter.RoomID)
	}

	var total int64
	query.Count(&total)

	var reservations []models.Reservation
	if err := query.Preload("Room.Hotel").Preload("User").Preload("Payment").
		Limit(filter.Limit).
		Offset(utils.Offset(filter.Page, filter.Limit)).
		Find(&reservations).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Paginated(c, http.StatusOK, reservations, total, filter.Page, filter.Limit)
}

// GetReservation godoc
// GET /reservations/:id
func (h *ReservationHandler) GetReservation(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	role := c.MustGet("userRole").(string)
	id, _ := strconv.Atoi(c.Param("id"))

	var reservation models.Reservation
	if err := config.DB.Preload("Room.Hotel").Preload("User").Preload("Payment").
		First(&reservation, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "reservation not found")
		return
	}

	if role != "admin" && reservation.UserID != userID {
		utils.Fail(c, http.StatusForbidden, "access denied")
		return
	}

	utils.OK(c, http.StatusOK, "", reservation)
}

// CreateReservation godoc
// POST /reservations
func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req models.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate dates
	if req.CheckOut.Before(req.CheckIn) || req.CheckOut.Equal(req.CheckIn) {
		utils.Fail(c, http.StatusBadRequest, "check_out must be after check_in")
		return
	}
	if req.CheckIn.Before(time.Now()) {
		utils.Fail(c, http.StatusBadRequest, "check_in must be in the future")
		return
	}

	var room models.Room
	if err := config.DB.First(&room, req.RoomID).Error; err != nil {
		utils.Fail(c, http.StatusBadRequest, "room not found")
		return
	}

	if !room.IsAvailable {
		utils.Fail(c, http.StatusConflict, "room is not available")
		return
	}

	// Check for overlapping reservations
	var count int64
	config.DB.Model(&models.Reservation{}).Where(
		"room_id = ? AND status IN ('pending','confirmed') AND check_in < ? AND check_out > ?",
		req.RoomID, req.CheckOut, req.CheckIn,
	).Count(&count)

	if count > 0 {
		utils.Fail(c, http.StatusConflict, "room is already booked for those dates")
		return
	}

	// Calculate price
	days := int(req.CheckOut.Sub(req.CheckIn).Hours() / 24)
	totalPrice := float64(days) * room.Price

	reservation := models.Reservation{
		UserID:     userID,
		RoomID:     req.RoomID,
		CheckIn:    req.CheckIn,
		CheckOut:   req.CheckOut,
		Guests:     req.Guests,
		TotalPrice: totalPrice,
		Status:     models.StatusPending,
		Notes:      req.Notes,
	}

	if err := config.DB.Create(&reservation).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	config.DB.Preload("Room.Hotel").Preload("User").First(&reservation, reservation.ID)
	utils.OK(c, http.StatusCreated, "reservation created", reservation)
}

// UpdateReservation godoc
// PUT /reservations/:id
func (h *ReservationHandler) UpdateReservation(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	role := c.MustGet("userRole").(string)
	id, _ := strconv.Atoi(c.Param("id"))

	var reservation models.Reservation
	if err := config.DB.First(&reservation, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "reservation not found")
		return
	}

	if role != "admin" && reservation.UserID != userID {
		utils.Fail(c, http.StatusForbidden, "access denied")
		return
	}

	var req models.UpdateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.CheckIn != nil {
		updates["check_in"] = req.CheckIn
	}
	if req.CheckOut != nil {
		updates["check_out"] = req.CheckOut
	}
	if req.Guests > 0 {
		updates["guests"] = req.Guests
	}
	if req.Notes != "" {
		updates["notes"] = req.Notes
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	config.DB.Model(&reservation).Updates(updates)
	config.DB.Preload("Room.Hotel").Preload("User").Preload("Payment").First(&reservation, reservation.ID)
	utils.OK(c, http.StatusOK, "reservation updated", reservation)
}

// CancelReservation godoc
// DELETE /reservations/:id
func (h *ReservationHandler) CancelReservation(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	role := c.MustGet("userRole").(string)
	id, _ := strconv.Atoi(c.Param("id"))

	var reservation models.Reservation
	if err := config.DB.First(&reservation, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "reservation not found")
		return
	}

	if role != "admin" && reservation.UserID != userID {
		utils.Fail(c, http.StatusForbidden, "access denied")
		return
	}

	config.DB.Model(&reservation).Update("status", models.StatusCancelled)
	utils.OK(c, http.StatusOK, "reservation cancelled", nil)
}
