package workout

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	CreateExercise(w http.ResponseWriter, r *http.Request)
	GetExerciseByID(w http.ResponseWriter, r *http.Request)
	ListExercises(w http.ResponseWriter, r *http.Request)
	UpdateExercise(w http.ResponseWriter, r *http.Request)
	DeleteExercise(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc Service
}

func NewHandler(svc Service) Handler {
	return &handler{svc: svc}
}

func (h *handler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	var exercise Exercise
	if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.svc.CreateExercise(r.Context(), exercise); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handler) GetExerciseByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	exercise, err := h.svc.GetExerciseByID(r.Context(), idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercise)
}

func (h *handler) ListExercises(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.svc.ListExercises(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercises)
}

func (h *handler) UpdateExercise(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	var exercise Exercise
	if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	exercise.ID = idStr

	if err := h.svc.UpdateExercise(r.Context(), exercise); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	if err := h.svc.DeleteExercise(r.Context(), idStr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
