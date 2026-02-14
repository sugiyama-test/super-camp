package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Time     string `json:"time"`
}

func HealthCheck(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		dbStatus := "up"
		if err := pool.Ping(ctx); err != nil {
			dbStatus = "down"
		}

		status := "ok"
		httpCode := http.StatusOK
		if dbStatus == "down" {
			status = "degraded"
			httpCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpCode)
		json.NewEncoder(w).Encode(HealthResponse{
			Status:   status,
			Database: dbStatus,
			Time:     time.Now().UTC().Format(time.RFC3339),
		})
	}
}
