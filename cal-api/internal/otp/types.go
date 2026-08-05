package otp

import "time"

type Otp struct {
	Email     string
	Code      string
	ExpiresAt time.Time
}

type OtpLogic interface {
	Create(email string) (string, error)
	VerifyAndConsume(email, otp string) (bool, error)
}
