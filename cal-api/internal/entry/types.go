package entry

import "time"

type EntryDB struct {
	Id     string
	UserId string
}

type Entry struct {
	id     string
	userId string
	value  EntryValue
	start  time.Time
	end    time.Time
}

type EntryLogic interface {
	Create(userId string, start, end time.Time) (string, error)
	Read(id string) (Entry, error)
	Update(id string, value EntryValue) (string, error)
	Delete(id string) (string, error)
	List(userId string, start, end time.Time) ([]Entry, error)
}
