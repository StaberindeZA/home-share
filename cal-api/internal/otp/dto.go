package otp

type RequestOtpDTO struct {
	Email string `json:"email"`
}

type VerifyOtpDTO struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

type OtpTokenDTO struct {
	Token string `json:"token"`
}
