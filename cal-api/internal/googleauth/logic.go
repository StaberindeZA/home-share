package googleauth

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"unicode/utf8"

	"cal-api/internal/user"

	"google.golang.org/api/idtoken"
)

var ErrClientIDNotFound = errors.New("google client id env var was not found")

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

func (gal SimpleGALogic) ValidateIDToken(idToken string) (string, error) {
	clientID, err := getGoogleClientID()
	if err != nil {
		return "", err
	}
	payload, err := idtoken.Validate(context.Background(), idToken, clientID)
	if err != nil {
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

	return email, nil
}

func NewSimpleGALogic(db *sql.DB, userLogic user.UserLogic) SimpleGALogic {
	return SimpleGALogic{
		db:        db,
		userLogic: userLogic,
	}
}
