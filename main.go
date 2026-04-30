package main

import (
	"log"

	"hotel-booking/config"
	"hotel-booking/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	gin.SetMode(cfg.GinMode)

	// Connect database
	db, err := config.ConnectDB(&cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Run migrations
	if err := config.AutoMigrate(db); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	log.Println("✅ Migrations applied")

	// Setup router
	r := handlers.SetupRouter(cfg)

	log.Printf("🚀 Server running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
