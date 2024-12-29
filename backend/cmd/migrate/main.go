package main

import (
	"backend/internal/db"
	"log"
)

func main() {
	err := db.RunMigrations()
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	} else {
		log.Println("Migrations completed successfully")
	}
}
