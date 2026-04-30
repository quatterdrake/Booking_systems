package handlers

import (
	"net/http"

	"hotel-booking/config"
	"hotel-booking/models"
	"hotel-booking/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

// Register godoc
// POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check duplicate email
	var existing models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		utils.Fail(c, http.StatusConflict, "email already registered")
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     models.RoleGuest,
	}

	if err := user.HashPassword(); err != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to hash password")
		return
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role),
		h.cfg.JWT.Secret, h.cfg.JWT.ExpirationHours)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	utils.OK(c, http.StatusCreated, "registered successfully", models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// Login godoc
// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.Fail(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if !user.CheckPassword(req.Password) {
		utils.Fail(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role),
		h.cfg.JWT.Secret, h.cfg.JWT.ExpirationHours)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	utils.OK(c, http.StatusOK, "login successful", models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// Me godoc
// GET /auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "user not found")
		return
	}

	utils.OK(c, http.StatusOK, "", user)
}
