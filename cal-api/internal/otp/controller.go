// Package otp provides HTTP handlers for OTP requests and verification
package otp

import (
	"encoding/json"
	"log"
	"net/http"

	"cal-api/internal/email"
	"cal-api/internal/utils"
)

type OtpController struct {
	logic       OtpLogic
	emailSender email.GomailSender
}

func (c OtpController) RequestOtp(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var requestOtp RequestOtpDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestOtp)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	otp, err := c.logic.Create(requestOtp.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err := c.emailSender.SendOtp(requestOtp.Email, otp); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c OtpController) VerifyOtp(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var verifyOtp VerifyOtpDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&verifyOtp)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	isValid, err := c.logic.VerifyAndConsume(verifyOtp.Email, verifyOtp.Otp)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !isValid {
		http.Error(w, "Invalid OTP code", http.StatusBadRequest)
		return
	}

	// Issue custom backend token
	tokenString, err := utils.GenerateBackendJWT(verifyOtp.Email)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := OtpTokenDTO{
		Token: tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func NewOtpController(logic OtpLogic, emailSender email.GomailSender) OtpController {
	return OtpController{
		logic:       logic,
		emailSender: emailSender,
	}
}
