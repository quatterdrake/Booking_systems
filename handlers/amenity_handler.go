package handlers

import (
	"net/http"

	"hotel-booking/config"
	"hotel-booking/models"
	"hotel-booking/utils"

	"github.com/gin-gonic/gin"
)

type AmenityHandler struct{}

func NewAmenityHandler() *AmenityHandler { return &AmenityHandler{} }

// GetAmenities godoc
// GET /amenities
func (h *AmenityHandler) GetAmenities(c *gin.Context) {
	var amenities []models.Amenity
	if err := config.DB.Find(&amenities).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.OK(c, http.StatusOK, "", amenities)
}

// CreateAmenity godoc
// POST /amenities
func (h *AmenityHandler) CreateAmenity(c *gin.Context) {
	var req models.CreateAmenityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	amenity := models.Amenity{Name: req.Name, Icon: req.Icon}
	if err := config.DB.Create(&amenity).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.OK(c, http.StatusCreated, "amenity created", amenity)
}
