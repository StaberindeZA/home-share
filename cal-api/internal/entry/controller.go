package entry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type EntryController struct {
	logic EntryLogic
}

func (c EntryController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var createEntry CreateEntryDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&createEntry)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	fmt.Println("here are the values from create entry")
	fmt.Println(createEntry)

	message, err := c.logic.Create(createEntry.UserId, createEntry.Start, createEntry.End)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(message))
}

func (c EntryController) Read(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	entry, err := c.logic.Read(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	data := EntryDTO{
		Id:     entry.id,
		UserId: entry.userId,
		Start:  entry.start.String(),
		End:    entry.end.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (c EntryController) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	message, err := c.logic.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(message))
}

func (c EntryController) List(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "Missing userId query param", http.StatusBadRequest)
	}
	start := r.URL.Query().Get("start")
	var startTime time.Time
	var err error
	if start == "" {
		startTime = time.Now()
	} else {
		startTime, err = time.Parse(time.RFC3339, start)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}
	end := r.URL.Query().Get("end")
	var endTime time.Time
	if end == "" {
		endTime = time.Now()
	} else {
		endTime, err = time.Parse(time.RFC3339, end)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}
	entries, err := c.logic.List(userId, startTime, endTime)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	data := make([]EntryDTO, 0, 1)
	for _, entry := range entries {
		data = append(data, EntryDTO{
			Id:     entry.id,
			UserId: entry.userId,
			Start:  entry.start.String(),
			End:    entry.end.String(),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func NewEntryController(logic EntryLogic) EntryController {
	return EntryController{
		logic: logic,
	}
}
