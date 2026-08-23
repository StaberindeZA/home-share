// Package healthcheck provides APIs to query health of the service
package healthcheck

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"cal-api/internal/utils"
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
	data := HealthResponse{
		Status:    "UP",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	utils.SendPayloadJSON(w, data)
}

func (c HealthcheckController) Ready(w http.ResponseWriter, r *http.Request) {
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

	utils.SendPayloadJSON(w, response)
}
