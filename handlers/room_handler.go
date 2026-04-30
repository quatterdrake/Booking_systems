package handlers

import (
	"net/http"
	"strconv"

	"hotel-booking/config"
	"hotel-booking/models"
	"hotel-booking/utils"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct{}

func NewRoomHandler() *RoomHandler { return &RoomHandler{} }

// GetRooms godoc
// GET /rooms
func (h *RoomHandler) GetRooms(c *gin.Context) {
	var filter models.RoomFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	query := config.DB.Model(&models.Room{})

	if filter.HotelID > 0 {
		query = query.Where("hotel_id = ?", filter.HotelID)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}
	if filter.Capacity > 0 {
		query = query.Where("capacity >= ?", filter.Capacity)
	}
	if filter.IsAvailable != nil {
		query = query.Where("is_available = ?", *filter.IsAvailable)
	}

	var total int64
	query.Count(&total)

	var rooms []models.Room
	if err := query.Preload("Hotel").Preload("Amenities").
		Limit(filter.Limit).
		Offset(utils.Offset(filter.Page, filter.Limit)).
		Find(&rooms).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Paginated(c, http.StatusOK, rooms, total, filter.Page, filter.Limit)
}

// GetRoom godoc
// GET /rooms/:id
func (h *RoomHandler) GetRoom(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var room models.Room
	if err := config.DB.Preload("Hotel").Preload("Amenities").First(&room, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "room not found")
		return
	}

	utils.OK(c, http.StatusOK, "", room)
}

// CreateRoom godoc
// POST /rooms
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req models.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	// Verify hotel exists
	var hotel models.Hotel
	if err := config.DB.First(&hotel, req.HotelID).Error; err != nil {
		utils.Fail(c, http.StatusBadRequest, "hotel not found")
		return
	}

	room := models.Room{
		Number:      req.Number,
		Type:        req.Type,
		Description: req.Description,
		Price:       req.Price,
		Capacity:    req.Capacity,
		HotelID:     req.HotelID,
		IsAvailable: true,
	}

	if err := config.DB.Create(&room).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Associate amenities
	if len(req.AmenityIDs) > 0 {
		var amenities []models.Amenity
		config.DB.Find(&amenities, req.AmenityIDs)
		config.DB.Model(&room).Association("Amenities").Replace(amenities)
	}

	config.DB.Preload("Hotel").Preload("Amenities").First(&room, room.ID)
	utils.OK(c, http.StatusCreated, "room created", room)
}

// UpdateRoom godoc
// PUT /rooms/:id
func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var room models.Room
	if err := config.DB.First(&room, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "room not found")
		return
	}

	var req models.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Number != "" {
		updates["number"] = req.Number
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if req.Capacity > 0 {
		updates["capacity"] = req.Capacity
	}
	if req.IsAvailable != nil {
		updates["is_available"] = *req.IsAvailable
	}

	config.DB.Model(&room).Updates(updates)

	if len(req.AmenityIDs) > 0 {
		var amenities []models.Amenity
		config.DB.Find(&amenities, req.AmenityIDs)
		config.DB.Model(&room).Association("Amenities").Replace(amenities)
	}

	config.DB.Preload("Hotel").Preload("Amenities").First(&room, room.ID)
	utils.OK(c, http.StatusOK, "room updated", room)
}

// DeleteRoom godoc
// DELETE /rooms/:id
func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := config.DB.Delete(&models.Room{}, id).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.OK(c, http.StatusOK, "room deleted", nil)
}

// GetFavoriteRooms godoc
// GET /rooms/favorites
func (h *RoomHandler) GetFavoriteRooms(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	page := 1
	limit := 10
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}

	var total int64
	config.DB.Model(&models.FavoriteRoom{}).Where("user_id = ?", userID).Count(&total)

	var favorites []models.FavoriteRoom
	if err := config.DB.
		Preload("Room.Hotel").
		Preload("Room.Amenities").
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(utils.Offset(page, limit)).
		Find(&favorites).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Paginated(c, http.StatusOK, favorites, total, page, limit)
}

// AddFavoriteRoom godoc
// PUT /rooms/:id/favorites
func (h *RoomHandler) AddFavoriteRoom(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	roomID, _ := strconv.Atoi(c.Param("id"))

	var room models.Room
	if err := config.DB.First(&room, roomID).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "room not found")
		return
	}

	fav := models.FavoriteRoom{
		UserID: userID,
		RoomID: uint(roomID),
	}

	// Use FirstOrCreate to avoid duplicates
	result := config.DB.Where(fav).FirstOrCreate(&fav)
	if result.Error != nil {
		utils.Fail(c, http.StatusInternalServerError, result.Error.Error())
		return
	}

	utils.OK(c, http.StatusOK, "added to favorites", fav)
}

// RemoveFavoriteRoom godoc
// DELETE /rooms/:id/favorites
func (h *RoomHandler) RemoveFavoriteRoom(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	roomID, _ := strconv.Atoi(c.Param("id"))

	if err := config.DB.
		Where("user_id = ? AND room_id = ?", userID, roomID).
		Delete(&models.FavoriteRoom{}).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.OK(c, http.StatusOK, "removed from favorites", nil)
}
