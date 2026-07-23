package main

import (
	"database/sql"
	"fmt"
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

	_, err = db.Exec(`DROP TABLE IF EXISTS entries;`)
	if err != nil {
		log.Fatalf("Failed to DROP entries table: %v", err)
	}
	fmt.Println("Table entries dropped successfully.")

	_, err = db.Exec(`DROP TABLE IF EXISTS users;`)
	if err != nil {
		log.Fatalf("Failed to DROP users table: %v", err)
	}
	fmt.Println("Table users dropped successfully.")

	usersSchema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(usersSchema); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	entriesSchema := `
	CREATE TABLE IF NOT EXISTS entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER REFERENCES users(id),
		start DATETIME,
		end DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(entriesSchema); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	fmt.Println("Table verified successfully.")
}

