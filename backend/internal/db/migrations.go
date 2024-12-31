package db

import (
	"log"
	"os/exec"
)

func RunMigrations(databaseURL string) error {
	// Убедитесь, что путь и база данных указаны правильно
	cmd := exec.Command("migrate", "-path", "./migrations", "-database", databaseURL, "up")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Migration error: %s", string(output))
		return err
	}
	log.Println(string(output))
	return nil
}
