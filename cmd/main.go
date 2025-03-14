package main

import (
	"github.com/genryusaishigikuni/spy_cats/config"
	"github.com/genryusaishigikuni/spy_cats/pkg/database"
	"github.com/genryusaishigikuni/spy_cats/pkg/router"
	"log"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("Cannot connect DB: %v", err)
	}

	// Run database migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Setup router
	r := router.SetupRouter(db)

	// Start server
	if err := r.Run(cfg.ServerPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
