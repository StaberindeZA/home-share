// Package homemate
package homemate

import "time"

type HomeMate struct {
	id        int
	homeID    int
	mateID    int
	Role      Role
	createdAt time.Time
	updatedAt time.Time
}

type HomeMateLogic interface {
	Create(homeID, mateID int, role Role) error
	ReadForHomeAndMate(homeID, mateID int) (HomeMate, error)
	Delete(homeID, mateID int) error
	// ListForHomeID(homeID int) ([]HomeMate, error)
	// ListForMateID(mateID int) ([]HomeMate, error)
}
