package home

import "cal-api/internal/user"

type HomeLogic interface {
	RetrieveMates(homeSlug string) ([]user.User, error)
}
