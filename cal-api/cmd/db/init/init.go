package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	force := flag.Bool("force", false, "Force re-init of DB. WARNING!!! This drops the databases.")
	flag.Parse()

	if _, err := os.Stat("data/app.db"); err == nil {
		if !*force {
			log.Println("SQLite db already exists. Skipping init.")
			os.Exit(0)
		}
	}

	db, err := sql.Open("sqlite", "data/app.db?_time_format=sqlite")
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

	_, err = db.Exec(`DROP TABLE IF EXISTS homemates;`)
	if err != nil {
		log.Fatalf("Failed to DROP users table: %v", err)
	}
	fmt.Println("Table homemates dropped successfully.")

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

	_, err = db.Exec(`DROP TABLE IF EXISTS otps;`)
	if err != nil {
		log.Fatalf("Failed to DROP users table: %v", err)
	}
	fmt.Println("Table otps dropped successfully.")

	_, err = db.Exec(`DROP TABLE IF EXISTS homes;`)
	if err != nil {
		log.Fatalf("Failed to DROP users table: %v", err)
	}
	fmt.Println("Table homes dropped successfully.")

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
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		value INTEGER NOT NULL DEFAULT 0 CHECK(value IN(0, 1, 2, 3)),
		start DATETIME,
		end DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(entriesSchema); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	otpsSchema := `
	CREATE TABLE IF NOT EXISTS otps (
		email TEXT PRIMARY KEY,
		code TEXT NOT NULL,
		expires_at DATETIME
	);`
	if _, err := db.Exec(otpsSchema); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	homesSchema := `
	CREATE TABLE IF NOT EXISTS homes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		slug TEXT NOT NULL UNIQUE,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(homesSchema); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	homeMatesSchema := `
	CREATE TABLE IF NOT EXISTS homemates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		home_id INTEGER REFERENCES homes(id) ON DELETE CASCADE,
		mate_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		role INTEGER NOT NULL DEFAULT 0 CHECK(role IN(0, 1)),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(homeMatesSchema); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	fmt.Println("Table verified successfully.")
}
