package entry

import "time"

type EntryDB struct {
	Id     string
	UserId string
}

type Entry struct {
	id     int
	userId int
	value  EntryValue
	start  time.Time
	end    time.Time
}

type EntryLogic interface {
	Create(userId int, start, end time.Time) (string, error)
	Read(id int) (Entry, error)
	Update(id int, value EntryValue) (string, error)
	Delete(id int) (string, error)
	List(userId int, start, end time.Time) ([]Entry, error)
}
