// Package router holds all routes exposed by cal-api
package router

import (
	"database/sql"
	"net/http"

	"cal-api/internal/auth"
	"cal-api/internal/email"
	"cal-api/internal/entry"
	"cal-api/internal/googleauth"
	"cal-api/internal/home"
	"cal-api/internal/homemate"
	"cal-api/internal/middleware"
	"cal-api/internal/otp"
	"cal-api/internal/user"

	_ "modernc.org/sqlite"
)

func New(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	userLogic := user.NewLiteLogic(db)
	homeMateLogic := homemate.NewLiteLogic(db)
	homeLogic := home.NewLiteLogic(db, userLogic)

	emailMessages := email.NewGomailMessages()
	emailSender := email.NewGomailSender(emailMessages)

	otpLogic := otp.NewLiteOtp(db, userLogic)
	otpController := otp.NewOtpController(otpLogic, emailSender)
	mux.HandleFunc("POST /v1/otp", otpController.RequestOtp)
	mux.HandleFunc("POST /v1/otp/verify", otpController.VerifyOtp)

	googleAuthLogic := googleauth.NewSimpleGALogic(db, userLogic)
	googleAuthController := googleauth.NewGoogleAuthController(googleAuthLogic)
	mux.HandleFunc("POST /v1/auth/google", googleAuthController.VerifyIdToken)

	authMiddleware := auth.NewAuthMiddleware(userLogic)
	authController := auth.NewAuthController(userLogic, homeLogic, homeMateLogic)
	mux.Handle("GET /v1/auth/userinfo", authMiddleware.Protected(http.HandlerFunc(authController.UserInfo)))
	mux.Handle("POST /v1/auth/role/verify", authMiddleware.Protected(http.HandlerFunc(authController.VerifyRole)))

	entryLogic := entry.NewMVPLogic(db)
	entryController := entry.NewEntryController(entryLogic)
	mux.Handle("POST /v1/entry", authMiddleware.Protected(http.HandlerFunc(entryController.Create)))
	mux.Handle("/v1/entry/{id}", authMiddleware.Protected(http.HandlerFunc(entryController.Read)))
	mux.Handle("PUT /v1/entry/{id}", authMiddleware.Protected(http.HandlerFunc(entryController.Update)))
	mux.Handle("DELETE /v1/entry/{id}", authMiddleware.Protected(http.HandlerFunc(entryController.Delete)))
	mux.Handle("/v1/entry", authMiddleware.Protected(http.HandlerFunc(entryController.List)))

	userController := user.NewUserController(userLogic)
	mux.Handle("GET /v1/user", authMiddleware.Protected(http.HandlerFunc(userController.ReadLoggedInUser)))
	mux.Handle("PUT /v1/user", authMiddleware.Protected(http.HandlerFunc(userController.UpdateLoggedInUser)))

	homeController := home.NewHomeController(homeLogic, homeMateLogic, userLogic)
	mux.Handle("POST /v1/home", authMiddleware.Protected(http.HandlerFunc(homeController.Create)))
	mux.Handle("GET /v1/home/{slug}", authMiddleware.Protected(http.HandlerFunc(homeController.Read)))
	mux.Handle("DELETE /v1/home/{slug}", authMiddleware.Protected(http.HandlerFunc(homeController.Delete)))
	mux.Handle("GET /v1/homes", authMiddleware.Protected(http.HandlerFunc(homeController.List)))
	mux.Handle("POST /v1/home/{slug}/mate", authMiddleware.Protected(http.HandlerFunc(homeController.CreateHomeMate)))
	mux.Handle("DELETE /v1/home/{slug}/mate", authMiddleware.Protected(http.HandlerFunc(homeController.DeleteHomeMate)))
	mux.Handle("GET /v1/home/{slug}/mates", authMiddleware.Protected(http.HandlerFunc(homeController.ReadHomeMates)))

	return middleware.Logger(mux)
}
