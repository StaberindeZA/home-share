package auth

import (
	"cal-api/internal/user"
	"cal-api/internal/utils"

	"context"
	"fmt"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	userLogic user.UserLogic
}

func NewAuthMiddleware(userLogic user.UserLogic) AuthMiddleware {
	return AuthMiddleware{
		userLogic: userLogic,
	}
}

func (am AuthMiddleware) Protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or malformed authorization token", http.StatusUnauthorized)
			return
		}

		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		payload, err := utils.ParseAndValidateJWT(accessToken)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		email := fmt.Sprint(payload["email"])
		user, err := am.userLogic.FindByEmail(email)
		if err != nil {
			http.Error(w, "Invalid email in claim: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Token is valid! Add user information to request context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
