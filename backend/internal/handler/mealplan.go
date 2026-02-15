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

type MealPlanHandler struct {
	repo *repository.MealPlanRepository
}

func NewMealPlanHandler(repo *repository.MealPlanRepository) *MealPlanHandler {
	return &MealPlanHandler{repo: repo}
}

func (h *MealPlanHandler) List(w http.ResponseWriter, r *http.Request) {
	plans, err := h.repo.List(r.Context(), defaultUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if plans == nil {
		plans = []model.MealPlan{}
	}
	writeJSON(w, http.StatusOK, plans)
}

func (h *MealPlanHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	plan, err := h.repo.Get(r.Context(), id)
	if err == pgx.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, plan)
}

type createMealPlanRequest struct {
	Title    string `json:"title"`
	MealType string `json:"meal_type"`
	Servings int    `json:"servings"`
	Notes    string `json:"notes"`
}

func (h *MealPlanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createMealPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	if req.MealType == "" {
		req.MealType = "dinner"
	}
	if req.Servings <= 0 {
		req.Servings = 2
	}
	plan, err := h.repo.Create(r.Context(), defaultUserID, req.Title, req.MealType, req.Servings, req.Notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, plan)
}

type updateMealPlanRequest struct {
	Title    string `json:"title"`
	MealType string `json:"meal_type"`
	Servings int    `json:"servings"`
	Notes    string `json:"notes"`
}

func (h *MealPlanHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req updateMealPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Servings <= 0 {
		req.Servings = 2
	}
	plan, err := h.repo.Update(r.Context(), id, req.Title, req.MealType, req.Servings, req.Notes)
	if err == pgx.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, plan)
}

func (h *MealPlanHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
