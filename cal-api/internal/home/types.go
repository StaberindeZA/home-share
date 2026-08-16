// Package home, also commonly known as a workspace is a
// grouping of related users. Users can be members of
// many Homes
package home

import (
	"time"

	"cal-api/internal/homemate"
)

type Home struct {
	ID          int
	name        string
	slug        string
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

type ListHome struct {
	ID            int
	ListHomeMates []ListHomeMates
}

type ListHomeMates struct {
	ID    int
	Mates ListHomeMatesMates
}

type ListHomeMatesMates struct {
	ID    int
	Email string
	Name  string
	Role  homemate.Role
}

type HomeLogic interface {
	Create(name, description string) (string, error)
	Read(slug string) (Home, error)
	Delete(slug string) error
	List(mateID int) ([]Home, error)
	ReadMates(slug string) ([]ListHomeMatesMates, error)
}
