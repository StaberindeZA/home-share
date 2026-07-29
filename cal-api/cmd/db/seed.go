package main

import (
	"database/sql"
	"log"
	"time"

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
	userOneRes, err := db.Exec(insertUserQuery, "P", "p@example.com")
	if err != nil {
		log.Printf("Insert error: %v", err)
	}
	userTwoRes, err := db.Exec(insertUserQuery, "R", "r@example.com")
	if err != nil {
		log.Printf("Insert error: %v", err)
	}

	userOneId, err := userOneRes.LastInsertId()
	if err != nil {
		log.Printf("User One last InsertId error: %v", err)
	}
	userTwoId, err := userTwoRes.LastInsertId()
	if err != nil {
		log.Printf("User Two last InsertId error: %v", err)
	}

	insertEntriesQuery := `INSERT INTO entries (user_id, value, start, end) VALUES (?, ?, ?, ?);`
	if _, err = db.Exec(insertEntriesQuery, userOneId, 1, time.Now(), time.Now()); err != nil {
		log.Printf("User one entries insert error: %v", err)
	}
	if _, err = db.Exec(insertEntriesQuery, userTwoId, 2, time.Now(), time.Now()); err != nil {
		log.Printf("User two entries insert error: %v", err)
	}
}
