package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/sugiyamadaiki/super-camp/backend/internal/repository"
)

type ChecklistHandler struct {
	repo repository.ChecklistRepo
}

func NewChecklistHandler(repo repository.ChecklistRepo) *ChecklistHandler {
	return &ChecklistHandler{repo: repo}
}

const defaultUserID int64 = 1 // 認証実装まで仮ユーザー

func (h *ChecklistHandler) List(w http.ResponseWriter, r *http.Request) {
	checklists, err := h.repo.List(r.Context(), defaultUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if checklists == nil {
		checklists = []repository.ChecklistSummary{}
	}
	writeJSON(w, http.StatusOK, checklists)
}

func (h *ChecklistHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	checklist, err := h.repo.Get(r.Context(), id)
	if err == pgx.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, checklist)
}

type createChecklistRequest struct {
	Title string `json:"title"`
}

func (h *ChecklistHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createChecklistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	checklist, err := h.repo.Create(r.Context(), req.Title, defaultUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, checklist)
}

func (h *ChecklistHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

type addItemRequest struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

func (h *ChecklistHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	checklistID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req addItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}
	item, err := h.repo.AddItem(r.Context(), checklistID, req.Name, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

type updateItemRequest struct {
	Name      string `json:"name"`
	IsChecked bool   `json:"is_checked"`
	Quantity  int    `json:"quantity"`
}

func (h *ChecklistHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.ParseInt(chi.URLParam(r, "itemID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}
	var req updateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}
	item, err := h.repo.UpdateItem(r.Context(), itemID, req.Name, req.IsChecked, req.Quantity)
	if err == pgx.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ChecklistHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.ParseInt(chi.URLParam(r, "itemID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}
	if err := h.repo.DeleteItem(r.Context(), itemID); err == pgx.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
