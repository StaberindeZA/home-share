package entry

import "time"

type CreateEntryDTO struct {
	UserId string    `json:"userId"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}

type UpdateEntryDTO struct {
	Value int `json:"value"`
}

type EntryDTO struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
	Value  int    `json:"value"`
	Start  string `json:"start"`
	End    string `json:"end"`
}
