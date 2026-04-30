package handlers

import (
	"net/http"
	"strconv"

	"hotel-booking/config"
	"hotel-booking/models"
	"hotel-booking/utils"

	"github.com/gin-gonic/gin"
)

type HotelHandler struct{}

func NewHotelHandler() *HotelHandler { return &HotelHandler{} }

// GetHotels godoc
// GET /hotels
func (h *HotelHandler) GetHotels(c *gin.Context) {
	var filter models.HotelFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	query := config.DB.Model(&models.Hotel{})

	if filter.City != "" {
		query = query.Where("city ILIKE ?", "%"+filter.City+"%")
	}
	if filter.Country != "" {
		query = query.Where("country ILIKE ?", "%"+filter.Country+"%")
	}
	if filter.Stars > 0 {
		query = query.Where("stars = ?", filter.Stars)
	}

	var total int64
	query.Count(&total)

	var hotels []models.Hotel
	if err := query.Preload("Rooms").
		Limit(filter.Limit).
		Offset(utils.Offset(filter.Page, filter.Limit)).
		Find(&hotels).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Paginated(c, http.StatusOK, hotels, total, filter.Page, filter.Limit)
}

// GetHotel godoc
// GET /hotels/:id
func (h *HotelHandler) GetHotel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var hotel models.Hotel
	if err := config.DB.Preload("Rooms").First(&hotel, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "hotel not found")
		return
	}

	utils.OK(c, http.StatusOK, "", hotel)
}

// CreateHotel godoc
// POST /hotels
func (h *HotelHandler) CreateHotel(c *gin.Context) {
	var req models.CreateHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	hotel := models.Hotel{
		Name:        req.Name,
		Address:     req.Address,
		City:        req.City,
		Country:     req.Country,
		Stars:       req.Stars,
		Description: req.Description,
		Phone:       req.Phone,
		Email:       req.Email,
	}

	if err := config.DB.Create(&hotel).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.OK(c, http.StatusCreated, "hotel created", hotel)
}

// UpdateHotel godoc
// PUT /hotels/:id
func (h *HotelHandler) UpdateHotel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var hotel models.Hotel
	if err := config.DB.First(&hotel, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "hotel not found")
		return
	}

	var req models.UpdateHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	config.DB.Model(&hotel).Updates(req)
	utils.OK(c, http.StatusOK, "hotel updated", hotel)
}

// DeleteHotel godoc
// DELETE /hotels/:id
func (h *HotelHandler) DeleteHotel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := config.DB.Delete(&models.Hotel{}, id).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.OK(c, http.StatusOK, "hotel deleted", nil)
}
