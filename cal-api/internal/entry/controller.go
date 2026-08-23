package entry

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"cal-api/internal/user"
	"cal-api/internal/utils"
)

type EntryController struct {
	logic EntryLogic
}

func NewEntryController(logic EntryLogic) EntryController {
	return EntryController{
		logic: logic,
	}
}

func (c EntryController) Create(w http.ResponseWriter, r *http.Request) {
	ok := utils.CheckContentTypeJSON(w, r)
	if !ok {
		return
	}
	user, ok := user.RetrieveUserFromContext(w, r)
	if !ok {
		return
	}

	createEntry, ok := utils.DecodeBodyJSON[CreateEntryDTO](w, r)
	if !ok {
		return
	}

	message, err := c.logic.Create(user.Id, createEntry.Start, createEntry.End)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(message))
}

func (c EntryController) Read(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	// TODO - With this function a double read on entity occurs
	if err := c.VerifyUserForEntity(id, r); err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

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
	utils.SendPayloadJSON(w, data)
}

func (c EntryController) Update(w http.ResponseWriter, r *http.Request) {
	ok := utils.CheckContentTypeJSON(w, r)
	if !ok {
		return
	}

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	if err := c.VerifyUserForEntity(id, r); err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	updateEntry, ok := utils.DecodeBodyJSON[UpdateEntryDTO](w, r)
	if !ok {
		return
	}

	message, err := c.logic.Update(id, EntryValue(updateEntry.Value))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(message))
}

func (c EntryController) Delete(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	if err := c.VerifyUserForEntity(id, r); err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	message, err := c.logic.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(message))
}

func (c EntryController) List(w http.ResponseWriter, r *http.Request) {
	userIDString := r.URL.Query().Get("userId")
	if userIDString == "" {
		http.Error(w, "Missing userId query param", http.StatusBadRequest)
	}
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	var startTime time.Time
	start := r.URL.Query().Get("start")
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
	entries, err := c.logic.List(userID, startTime, endTime)
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
			Value:  int(entry.value),
			Start:  entry.start.String(),
			End:    entry.end.String(),
		})
	}

	utils.SendPayloadJSON(w, data)
}

func (c EntryController) VerifyUserForEntity(entryID int, r *http.Request) error {
	user, ok := r.Context().Value("user").(user.User)
	if !ok {
		return errors.New("user not in context")
	}

	entry, err := c.logic.Read(entryID)
	if err != nil {
		return err
	}

	if entry.userId != user.Id {
		return errors.New("user does not match entity user")
	}

	return nil
}
