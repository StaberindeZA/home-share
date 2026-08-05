package entry

import "time"

type CreateEntryDTO struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type UpdateEntryDTO struct {
	Value int `json:"value"`
}

type EntryDTO struct {
	Id     int    `json:"id"`
	UserId int    `json:"userId"`
	Value  int    `json:"value"`
	Start  string `json:"start"`
	End    string `json:"end"`
}
