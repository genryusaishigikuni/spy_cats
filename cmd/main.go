package main

import (
	"log"

	"github.com/genryusaishigikuni/spy_cats/config"
	"github.com/genryusaishigikuni/spy_cats/pkg/database"
	"github.com/genryusaishigikuni/spy_cats/pkg/router"
)

func main() {
	// 1. Here I Load config
	cfg := config.Load()

	// 2. Here I establish the connection with DB
	db, err := database.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("Cannot connect DB: %v", err)
	}

	// Now I run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 4. Here I set up router
	r := router.SetupRouter(db)

	// 5. Finally, I start the server
	if err := r.Run(cfg.ServerPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
