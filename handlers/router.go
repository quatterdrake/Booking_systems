package handlers

import (
	"hotel-booking/config"
	"hotel-booking/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "hotel-booking"})
	})

	// Handlers
	authH := NewAuthHandler(cfg)
	hotelH := NewHotelHandler()
	roomH := NewRoomHandler()
	reservationH := NewReservationHandler()
	paymentH := NewPaymentHandler()
	amenityH := NewAmenityHandler()

	auth := middleware.AuthMiddleware(cfg.JWT.Secret)
	adminOnly := middleware.AdminOnly()

	// ── Auth routes ──────────────────────────────
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authH.Register)
		authGroup.POST("/login", authH.Login)
		authGroup.GET("/me", auth, authH.Me)
	}

	// ── Hotel routes ─────────────────────────────
	hotels := r.Group("/hotels")
	{
		hotels.GET("", hotelH.GetHotels)
		hotels.GET("/:id", hotelH.GetHotel)
		hotels.POST("", auth, adminOnly, hotelH.CreateHotel)
		hotels.PUT("/:id", auth, adminOnly, hotelH.UpdateHotel)
		hotels.DELETE("/:id", auth, adminOnly, hotelH.DeleteHotel)
	}

	// ── Room routes ──────────────────────────────
	rooms := r.Group("/rooms")
	{
		rooms.GET("", roomH.GetRooms)
		rooms.GET("/favorites", auth, roomH.GetFavoriteRooms)  // Practicum 6
		rooms.GET("/:id", roomH.GetRoom)
		rooms.POST("", auth, adminOnly, roomH.CreateRoom)
		rooms.PUT("/:id", auth, adminOnly, roomH.UpdateRoom)
		rooms.DELETE("/:id", auth, adminOnly, roomH.DeleteRoom)
		rooms.PUT("/:id/favorites", auth, roomH.AddFavoriteRoom)    // Practicum 6
		rooms.DELETE("/:id/favorites", auth, roomH.RemoveFavoriteRoom) // Practicum 6
	}

	// ── Reservation routes ───────────────────────
	reservations := r.Group("/reservations", auth)
	{
		reservations.GET("", reservationH.GetReservations)
		reservations.GET("/:id", reservationH.GetReservation)
		reservations.POST("", reservationH.CreateReservation)
		reservations.PUT("/:id", reservationH.UpdateReservation)
		reservations.DELETE("/:id", reservationH.CancelReservation)
	}

	// ── Payment routes ───────────────────────────
	payments := r.Group("/payments", auth)
	{
		payments.POST("", paymentH.CreatePayment)
		payments.GET("/:id", paymentH.GetPayment)
	}

	// ── Amenity routes ───────────────────────────
	amenities := r.Group("/amenities")
	{
		amenities.GET("", amenityH.GetAmenities)
		amenities.POST("", auth, adminOnly, amenityH.CreateAmenity)
	}

	return r
}
