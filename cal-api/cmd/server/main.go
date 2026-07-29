package main

import (
	"cal-api/internal/router"
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "assets/app.db?_time_format=sqlite")
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
