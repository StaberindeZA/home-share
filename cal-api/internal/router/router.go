// Package router holds all routes exposed by cal-api
package router

import (
	"database/sql"
	"net/http"

	"cal-api/internal/auth"
	"cal-api/internal/email"
	"cal-api/internal/entry"
	"cal-api/internal/googleauth"
	"cal-api/internal/healthcheck"
	"cal-api/internal/home"
	"cal-api/internal/homemate"
	"cal-api/internal/middleware"
	"cal-api/internal/otp"
	"cal-api/internal/user"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "modernc.org/sqlite"
)

func New(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	promReg := prometheus.NewRegistry()
	promReg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector((collectors.ProcessCollectorOpts{})),
	)

	httpMetrics := middleware.NewHTTPMetrics(promReg)

	userLogic := user.NewLiteLogic(db)
	homeMateLogic := homemate.NewLiteLogic(db)
	entryLogic := entry.NewMVPLogic(db)
	homeLogic := home.NewLiteLogic(db, userLogic)
	otpLogic := otp.NewLiteOtp(db, userLogic)
	googleAuthLogic := googleauth.NewSimpleGALogic(db, userLogic)
	authMiddleware := auth.NewAuthMiddleware(userLogic)

	emailMessages := email.NewGomailMessages()
	emailSender := email.NewGomailSender(emailMessages)

	healthcheckController := healthcheck.NewHealthcheckController(db)
	mux.HandleFunc("/livez", healthcheckController.Live)
	mux.HandleFunc("/readyz", healthcheckController.Ready)

	otpController := otp.NewOtpController(otpLogic, emailSender)
	mux.HandleFunc("POST /v1/otp", otpController.RequestOtp)
	mux.HandleFunc("POST /v1/otp/verify", otpController.VerifyOtp)

	googleAuthController := googleauth.NewGoogleAuthController(googleAuthLogic)
	mux.HandleFunc("POST /v1/auth/google", googleAuthController.VerifyIDToken)

	authController := auth.NewAuthController(userLogic, homeLogic, homeMateLogic)
	mux.Handle("GET /v1/auth/userinfo", authMiddleware.Protected(http.HandlerFunc(authController.UserInfo)))
	mux.Handle("POST /v1/auth/role/verify", authMiddleware.Protected(http.HandlerFunc(authController.VerifyRole)))

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

	mux.Handle("/metrics", promhttp.HandlerFor(promReg, promhttp.HandlerOpts{}))

	return middleware.Logger(httpMetrics.Middleware(mux))
}
