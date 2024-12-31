package main

import (
	"backend/internal/config"
	"backend/internal/db"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	err := db.RunMigrations(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	} else {
		log.Println("Migrations completed successfully")
	}
}
