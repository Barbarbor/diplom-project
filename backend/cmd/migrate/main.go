package main

import (
	"backend/internal/config"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

const migrationsTable = "schema_migrations"

func ensureMigrationsTable(db *sql.DB) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id SERIAL PRIMARY KEY,
		version VARCHAR(255) UNIQUE NOT NULL,
		applied_at TIMESTAMP DEFAULT NOW()
	);`, migrationsTable)
	_, err := db.Exec(query)
	return err
}

func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	applied := make(map[string]bool)
	rows, err := db.Query(fmt.Sprintf("SELECT version FROM %s", migrationsTable))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}
	return applied, nil
}

func applyMigration(db *sql.DB, version, sqlStmt string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(sqlStmt); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute migration %s: %w", version, err)
	}
	_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (version) VALUES ($1)", migrationsTable), version)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record migration %s: %w", version, err)
	}
	return tx.Commit()
}

func runMigrations(db *sql.DB, migrationsDir string) error {
	// Убедимся, что таблица для миграций существует
	if err := ensureMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to ensure migrations table: %w", err)
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	// Сортировка файлов по имени (предполагается, что имя начинается с номера миграции)
	sort.Strings(migrationFiles)

	for _, filename := range migrationFiles {
		version := strings.TrimSuffix(filename, ".up.sql")
		if applied[version] {
			log.Printf("Migration %s already applied, skipping.", version)
			continue
		}
		path := filepath.Join(migrationsDir, filename)
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}
		sqlStmt := string(content)
		log.Printf("Applying migration %s", version)
		if err := applyMigration(db, version, sqlStmt); err != nil {
			return err
		}
		log.Printf("Migration %s applied successfully.", version)
	}

	return nil
}

func main() {
	cfg := config.LoadConfig()
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db, "./migrations"); err != nil {
		log.Fatalf("Migration failed: %v", err)
	} else {
		log.Println("Migrations completed successfully")
	}
}
