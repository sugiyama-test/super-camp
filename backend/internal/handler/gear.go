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

type GearHandler struct {
	repo repository.GearRepo
}

func NewGearHandler(repo repository.GearRepo) *GearHandler {
	return &GearHandler{repo: repo}
}

func (h *GearHandler) List(w http.ResponseWriter, r *http.Request) {
	gears, err := h.repo.List(r.Context(), defaultUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if gears == nil {
		gears = []model.Gear{}
	}
	writeJSON(w, http.StatusOK, gears)
}

func (h *GearHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	gear, err := h.repo.Get(r.Context(), id)
	if err == pgx.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, gear)
}

type createGearRequest struct {
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Brand       string   `json:"brand"`
	WeightGrams *float64 `json:"weight_grams"`
	Notes       string   `json:"notes"`
}

func (h *GearHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createGearRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	gear, err := h.repo.Create(r.Context(), defaultUserID, req.Name, req.Category, req.Brand, req.WeightGrams, req.Notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, gear)
}

type updateGearRequest struct {
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Brand       string   `json:"brand"`
	WeightGrams *float64 `json:"weight_grams"`
	Notes       string   `json:"notes"`
}

func (h *GearHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req updateGearRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	gear, err := h.repo.Update(r.Context(), id, req.Name, req.Category, req.Brand, req.WeightGrams, req.Notes)
	if err == pgx.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, gear)
}

func (h *GearHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
