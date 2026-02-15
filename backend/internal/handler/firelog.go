package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/sugiyamadaiki/super-camp/backend/internal/model"
	"github.com/sugiyamadaiki/super-camp/backend/internal/repository"
)

type FireLogHandler struct {
	repo *repository.FireLogRepository
}

func NewFireLogHandler(repo *repository.FireLogRepository) *FireLogHandler {
	return &FireLogHandler{repo: repo}
}

func (h *FireLogHandler) List(w http.ResponseWriter, r *http.Request) {
	logs, err := h.repo.List(r.Context(), defaultUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if logs == nil {
		logs = []model.FireLog{}
	}
	writeJSON(w, http.StatusOK, logs)
}

func (h *FireLogHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	log, err := h.repo.Get(r.Context(), id)
	if err == pgx.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, log)
}

type createFireLogRequest struct {
	Date            string   `json:"date"`
	Location        string   `json:"location"`
	WoodType        string   `json:"wood_type"`
	DurationMinutes int      `json:"duration_minutes"`
	Notes           string   `json:"notes"`
	Temperature     *float64 `json:"temperature"`
}

func (h *FireLogHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createFireLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Date == "" {
		http.Error(w, "date is required", http.StatusBadRequest)
		return
	}
	log, err := h.repo.Create(r.Context(), defaultUserID, req.Date, req.Location, req.WoodType, req.DurationMinutes, req.Notes, req.Temperature)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, log)
}

type updateFireLogRequest struct {
	Date            string   `json:"date"`
	Location        string   `json:"location"`
	WoodType        string   `json:"wood_type"`
	DurationMinutes int      `json:"duration_minutes"`
	Notes           string   `json:"notes"`
	Temperature     *float64 `json:"temperature"`
}

func (h *FireLogHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req updateFireLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	log, err := h.repo.Update(r.Context(), id, req.Date, req.Location, req.WoodType, req.DurationMinutes, req.Notes, req.Temperature)
	if err == pgx.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, log)
}

func (h *FireLogHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.repo.Delete(r.Context(), id); err == pgx.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
