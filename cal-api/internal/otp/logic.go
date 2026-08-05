package otp

import (
	"cal-api/internal/user"

	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"time"
	"unicode/utf8"
)

type LiteOtp struct {
	db        *sql.DB
	userLogic user.UserLogic
}

func (lo LiteOtp) Create(email string) (string, error) {
	existingOtp, err := lo.Read(email)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	if err != sql.ErrNoRows {
		if existingOtp.ExpiresAt.Before(time.Now()) {
			if err := lo.Delete(email); err != nil {
				return "", err
			}
		} else {
			return existingOtp.Code, nil
		}
	}

	insertOtpQuery := `INSERT INTO otps (email, code, expires_at) VALUES (?, ?, ?);`
	expiresAt := time.Now().Add(time.Minute * 5)

	max := big.NewInt(1000000)
	n, _ := rand.Int(rand.Reader, max)
	otp := fmt.Sprintf("%06d", n.Int64())

	if _, err := lo.db.Exec(insertOtpQuery, email, otp, expiresAt); err != nil {
		return "", err
	}

	return otp, nil
}

func (lo LiteOtp) Read(email string) (Otp, error) {
	var o Otp
	readQuery := `SELECT email, code, expires_at FROM otps WHERE email = ?;`
	err := lo.db.QueryRow(readQuery, email).Scan(&o.Email, &o.Code, &o.ExpiresAt)

	return o, err

}

func (lo LiteOtp) Delete(email string) error {
	removeOtpQuery := `DELETE FROM otps WHERE email = ?`
	_, err := lo.db.Exec(removeOtpQuery, email)
	return err
}

func (lo LiteOtp) VerifyAndConsume(email, otp string) (bool, error) {
	o, err := lo.Read(email)

	if err == sql.ErrNoRows {
		log.Println("OTP does not exist for email: %s", email)
	}

	if err != nil {
		return false, err
	}

	equalCode := o.Code == otp
	if equalCode == false {
		log.Println("Code does not match")
	}

	expiryTime := time.Now().Add(time.Minute * 5)
	unexpiredCode := expiryTime.After(o.ExpiresAt)
	if unexpiredCode == false {
		log.Println("Code has expired")
	}

	validOtp := equalCode && unexpiredCode

	if validOtp == true {
		emailFirstRune, _ := utf8.DecodeRuneInString(email)
		name := string(emailFirstRune)
		_, err := lo.userLogic.FindOrCreate(0, name, email)
		if err != nil {
			log.Fatal("Could not find or Create user")
		}

		if err := lo.Delete(email); err != nil {
			log.Fatal("Error deleting otp, %v", err)
			return false, err
		}
	}

	return validOtp, nil
}

func NewLiteOtp(db *sql.DB, userLogic user.UserLogic) LiteOtp {
	return LiteOtp{
		db:        db,
		userLogic: userLogic,
	}
}
