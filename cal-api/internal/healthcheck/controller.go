// Package healthcheck provides APIs to query health of the service
package healthcheck

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type HealthcheckController struct {
	db *sql.DB
}

func NewHealthcheckController(db *sql.DB) HealthcheckController {
	return HealthcheckController{
		db: db,
	}
}

func (c HealthcheckController) Live(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HealthResponse{
		Status:    "UP",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func (c HealthcheckController) Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	checks := make(map[string]string)
	isReady := true

	if err := c.db.PingContext(ctx); err != nil {
		checks["database"] = "DOWN: " + err.Error()
		isReady = false
	} else {
		checks["database"] = "UP"
	}

	response := HealthResponse{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	}

	if !isReady {
		response.Status = "DOWN"
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		response.Status = "UP"
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}
