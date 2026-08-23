// Package otp provides HTTP handlers for OTP requests and verification
package otp

import (
	"net/http"

	"cal-api/internal/email"
	"cal-api/internal/utils"
)

type OtpController struct {
	logic       OtpLogic
	emailSender email.GomailSender
}

func (c OtpController) RequestOtp(w http.ResponseWriter, r *http.Request) {
	if ok := utils.CheckContentTypeJSON(w, r); !ok {
		return
	}

	requestOtp, ok := utils.DecodeBodyJSON[RequestOtpDTO](w, r)
	if !ok {
		return
	}

	otp, err := c.logic.Create(requestOtp.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err := c.emailSender.SendOtp(requestOtp.Email, otp); err != nil {
		utils.RecordErrorServer(w, err, "Otp.RequestOtp.SendOtp")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c OtpController) VerifyOtp(w http.ResponseWriter, r *http.Request) {
	if ok := utils.CheckContentTypeJSON(w, r); !ok {
		return
	}

	verifyOtp, ok := utils.DecodeBodyJSON[VerifyOtpDTO](w, r)
	if !ok {
		return
	}

	isValid, err := c.logic.VerifyAndConsume(verifyOtp.Email, verifyOtp.Otp)
	if err != nil {
		utils.RecordErrorServer(w, err, "Otp.VerifyOtp.VerifyAndConsume")
		return
	}

	if !isValid {
		utils.RecordErrorServer(w, err, "Otp.VerifyOtp.isValid")
		return
	}

	tokenString, err := utils.GenerateBackendJWT(verifyOtp.Email)
	if err != nil {
		utils.RecordErrorServer(w, err, "Otp.VerifyOtp.GenerateBackendJWT")
		return
	}
	data := OtpTokenDTO{
		Token: tokenString,
	}

	utils.SendPayloadJSON(w, data)
}

func NewOtpController(logic OtpLogic, emailSender email.GomailSender) OtpController {
	return OtpController{
		logic:       logic,
		emailSender: emailSender,
	}
}
