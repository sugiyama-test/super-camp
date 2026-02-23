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

type CampsiteHandler struct {
	repo repository.CampsiteRepo
}

func NewCampsiteHandler(repo repository.CampsiteRepo) *CampsiteHandler {
	return &CampsiteHandler{repo: repo}
}

func (h *CampsiteHandler) List(w http.ResponseWriter, r *http.Request) {
	campsites, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if campsites == nil {
		campsites = []model.Campsite{}
	}
	writeJSON(w, http.StatusOK, campsites)
}

func (h *CampsiteHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	campsite, err := h.repo.Get(r.Context(), id)
	if err == pgx.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, campsite)
}

type createCampsiteRequest struct {
	Name      string   `json:"name"`
	Address   string   `json:"address"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Notes     string   `json:"notes"`
}

func (h *CampsiteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createCampsiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	campsite, err := h.repo.Create(r.Context(), req.Name, req.Address, req.Latitude, req.Longitude, req.Notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, campsite)
}

type updateCampsiteRequest struct {
	Name      string   `json:"name"`
	Address   string   `json:"address"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Notes     string   `json:"notes"`
}

func (h *CampsiteHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req updateCampsiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	campsite, err := h.repo.Update(r.Context(), id, req.Name, req.Address, req.Latitude, req.Longitude, req.Notes)
	if err == pgx.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, campsite)
}

func (h *CampsiteHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
