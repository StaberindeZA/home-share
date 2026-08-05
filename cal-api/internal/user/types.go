package user

import "time"

type User struct {
	Id        int
	Name      string
	Email     string
	createdAt time.Time
	updatedAt time.Time
}

type UserLogic interface {
	FindOrCreate(id int, name, email string) (User, error)
	FindByEmail(email string) (User, error)
}
