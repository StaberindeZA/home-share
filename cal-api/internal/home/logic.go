package home

import (
	"os"
	"strings"

	"cal-api/internal/user"
)

type StaticLogic struct {
	userLogic user.UserLogic
}

func NewStaticLogic(userLogic user.UserLogic) StaticLogic {
	return StaticLogic{
		userLogic: userLogic,
	}
}

func (l StaticLogic) RetrieveMates(homeSlug string) ([]user.User, error) {
	// TODO - Use the homeSlug to fetch users
	// Get list of emails from env var
	emailsString := os.Getenv("WORKSPACE_USER_EMAIL_STRING")
	emails := strings.Split(emailsString, ",")

	var mates []user.User
	var mateError error
	for _, email := range emails {
		mate, err := l.userLogic.FindByEmail(email)
		if err != nil {
			mateError = err
			break
		}
		mates = append(mates, mate)
	}

	if mateError != nil {
		var emptyMates []user.User
		return emptyMates, mateError
	}

	return mates, nil
}
