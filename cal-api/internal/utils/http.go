package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func CheckContentTypeJSON(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return false
	}
	return true
}

func DecodeBodyJSON[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var body T
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return body, false
	}

	return body, true
}

func SendPayloadJSON(w http.ResponseWriter, data any) bool {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		RecordErrorServer(w, err, "SendPayloadJSON")
		return false
	}

	return true
}

func RecordErrorServer(w http.ResponseWriter, err error, source string) {
	slog.Error("Internal Server Error", "error", err, "source", source)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
