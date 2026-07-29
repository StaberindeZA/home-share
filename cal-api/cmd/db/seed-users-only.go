package main

import (
	"database/sql"
	"log"

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

	// Assume that db/init.go has already been run and tables exist

	insertUserQuery := `INSERT INTO users (name, email) VALUES (?, ?);`
	_, err = db.Exec(insertUserQuery, "P", "p@example.com")
	if err != nil {
		log.Printf("Insert error: %v", err)
	}
	_, err = db.Exec(insertUserQuery, "R", "r@example.com")
	if err != nil {
		log.Printf("Insert error: %v", err)
	}
}
