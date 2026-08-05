package main

import (
	"cal-api/internal/otp"
	"cal-api/internal/user"
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
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

	liteLogic := user.NewLiteLogic(db)
	u, err := liteLogic.FindOrCreate(10, "Jeff", "j@example.com")
	fmt.Printf("Found user: %d, %s\n", u.Id, u.Name)

	liteOtp := otp.NewLiteOtp(db)
	_, err = liteOtp.Create("test@example.com", "12345")
	if err != nil {
		log.Println("Error during OTP create: %v", err)
	}
	validOtp, err := liteOtp.VerifyAndConsumeOtp("test@example.com", "12345")
	if err != nil {
		log.Println("Error during verify: %v", err)
	}

	log.Println("Was it a valid OTP: %s", validOtp)
}
