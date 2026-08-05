package googleauth

import (
	"cal-api/internal/user"

	"context"
	"database/sql"
	"errors"
	"os"
	"unicode/utf8"

	"google.golang.org/api/idtoken"
)

var ErrClientIDNotFound = errors.New("Google Client ID env var was not found.")

func getGoogleClientID() (string, error) {
	secret := os.Getenv("GOOGLE_CLIENT_ID")
	if secret == "" {
		return "", ErrClientIDNotFound
	}
	return secret, nil
}

type SimpleGALogic struct {
	db        *sql.DB
	userLogic user.UserLogic
}

func (gal SimpleGALogic) ValidateIdToken(idToken string) (string, error) {
	clientId, err := getGoogleClientID()
	if err != nil {
		return "", err
	}
	payload, err := idtoken.Validate(context.Background(), idToken, clientId)
	if err != nil {
		//http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return "", err
	}

	email := payload.Claims["email"].(string)

	// TODO - Move to util function
	emailFirstRune, _ := utf8.DecodeRuneInString(email)
	name := string(emailFirstRune)
	_, err = gal.userLogic.FindOrCreate(0, name, email)
	if err != nil {
		return "", err
	}

	// Issue custom backend token
	return email, nil
}

func NewSimpleGALogic(db *sql.DB, userLogic user.UserLogic) SimpleGALogic {
	return SimpleGALogic{
		db:        db,
		userLogic: userLogic,
	}
}
