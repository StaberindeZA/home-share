package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() string {
	secret := os.Getenv("AUTH_SECRET")
	if secret == "" {
		log.Fatalf("Missing env var: %s", "AUTH_SECRET")
	}
	return secret
}

func GenerateBackendJWT(email string) (string, error) {
	jwtSecret := []byte(getJWTSecret())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	})
	return token.SignedString(jwtSecret)
}

func ParseAndValidateJWT(tokenString string) (jwt.MapClaims, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(getJWTSecret()), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return jwt.MapClaims{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		log.Println(claims["email"], claims["nbf"])
		return claims, nil
	} else {
		log.Println(err)
		return jwt.MapClaims{}, err
	}
}
