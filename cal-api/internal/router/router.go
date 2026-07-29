package router

import (
	"cal-api/internal/entry"
	"cal-api/internal/middleware"
	"database/sql"
	"net/http"

	_ "modernc.org/sqlite"
)

func New(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	//entryLogic := entry.NewDummyLogic()
	entryLogic := entry.NewMVPLogic(db)
	entryController := entry.NewEntryController(entryLogic)
	mux.HandleFunc("POST /v1/entry", entryController.Create)
	mux.HandleFunc("/v1/entry/{id}", entryController.Read)
	mux.HandleFunc("PUT /v1/entry/{id}", entryController.Update)
	mux.HandleFunc("DELETE /v1/entry/{id}", entryController.Delete)
	mux.HandleFunc("/v1/entry", entryController.List)

	return middleware.Logger(mux)
}
