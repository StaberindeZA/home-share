package main

import (
	"database/sql"
	"embed"
	"errors"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite" // Pure-Go SQLite driver registration
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatalf("DB_HOST env var is required")
	}

	// 1. Open the SQLite database connection using the pure Go driver
	db, err := sql.Open("sqlite", "data/app.db")
	if err != nil {
		log.Fatalf("Could not open database: %v", err)
	}
	defer db.Close()

	// 2. Initialize the golang-migrate driver instance
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatalf("Could not create migration driver instance: %v", err)
	}

	sourceDriver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		log.Fatalf("Failed to create embedded source driver: %v", err)
	}

	// 3. Create a new migrator pointing to the filesystem directory
	m, err := migrate.NewWithInstance(
		"iofs",
		sourceDriver,
		"sqlite", // Database name identifier
		driver,
	)
	if err != nil {
		log.Fatalf("Migration initialization failed: %v", err)
	}

	// 4. Execute "Up" migrations
	log.Println("Running database migrations...")
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Database is already up to date.")
		} else {
			log.Fatalf("Migration failed: %v", err)
		}
	} else {
		log.Println("Migrations applied successfully!")
	}
}
