package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/sugiyamadaiki/super-camp/backend/internal/repository"
)

type LayoutHandler struct {
	repo *repository.LayoutRepository
}

func NewLayoutHandler(repo *repository.LayoutRepository) *LayoutHandler {
	return &LayoutHandler{repo: repo}
}

func (h *LayoutHandler) List(w http.ResponseWriter, r *http.Request) {
	layouts, err := h.repo.List(r.Context(), defaultUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if layouts == nil {
		layouts = []repository.LayoutSummary{}
	}
	writeJSON(w, http.StatusOK, layouts)
}

func (h *LayoutHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	layout, err := h.repo.Get(r.Context(), id)
	if err == pgx.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, layout)
}

type createLayoutRequest struct {
	Title string          `json:"title"`
	Data  json.RawMessage `json:"data"`
}

func (h *LayoutHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createLayoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	data := "{}"
	if req.Data != nil {
		data = string(req.Data)
	}
	layout, err := h.repo.Create(r.Context(), req.Title, data, defaultUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, layout)
}

type updateLayoutRequest struct {
	Title string          `json:"title"`
	Data  json.RawMessage `json:"data"`
}

func (h *LayoutHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req updateLayoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	data := "{}"
	if req.Data != nil {
		data = string(req.Data)
	}
	layout, err := h.repo.Update(r.Context(), id, req.Title, data)
	if err == pgx.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, layout)
}

func (h *LayoutHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
