package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"cal-api/internal/router"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	if _, err := os.Stat(dbHost); err != nil {
		log.Fatalf("SQL db does not exist at %s", dbHost)
		os.Exit(0)
	}

	db, err := sql.Open("sqlite", fmt.Sprintf("%s?_time_format=sqlite", dbHost))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(`PRAGMA journal_mode=WAL;`); err != nil {
		log.Fatalf("Failed to enable WAL: %v", err)
	}
	if _, err := db.Exec(`PRAGMA foreign_keys=ON;`); err != nil {
		log.Fatalf("Failed to enable Foreign Keys: %v", err)
	}

	appRouter := router.New(db)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: appRouter,
	}

	log.Println("Production server initializing on port 8080...")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Fatal bootstrap exit: %v", err)
	}
}
