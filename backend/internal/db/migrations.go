package db

import (
	"log"
	"os/exec"
)

func RunMigrations() error {
	// Убедитесь, что путь и база данных указаны правильно
	cmd := exec.Command("migrate", "-path", "./migrations", "-database", "postgres://postgres:12345@localhost:5432/surveydb?sslmode=disable", "up")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Migration error: %s", string(output))
		return err
	}
	log.Println(string(output))
	return nil
}
