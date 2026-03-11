package workout

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"nds-go-starter/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	CreateExercise(w http.ResponseWriter, r *http.Request)
	GetExerciseByID(w http.ResponseWriter, r *http.Request)
	ListExercises(w http.ResponseWriter, r *http.Request)
	UpdateExercise(w http.ResponseWriter, r *http.Request)
	DeleteExercise(w http.ResponseWriter, r *http.Request)

	CreateWorkoutSession(w http.ResponseWriter, r *http.Request)
	GetWorkoutSessionByID(w http.ResponseWriter, r *http.Request)
	ListWorkoutSessions(w http.ResponseWriter, r *http.Request)
	UpdateWorkoutSession(w http.ResponseWriter, r *http.Request)
	DeleteWorkoutSession(w http.ResponseWriter, r *http.Request)

	CreateWorkoutExercise(w http.ResponseWriter, r *http.Request)
	GetWorkoutExerciseByID(w http.ResponseWriter, r *http.Request)
	ListWorkoutExercisesBySession(w http.ResponseWriter, r *http.Request)
	UpdateWorkoutExercise(w http.ResponseWriter, r *http.Request)
	DeleteWorkoutExercise(w http.ResponseWriter, r *http.Request)

	CreateExerciseSet(w http.ResponseWriter, r *http.Request)
	GetExerciseSetByID(w http.ResponseWriter, r *http.Request)
	ListExerciseSets(w http.ResponseWriter, r *http.Request)
	UpdateExerciseSet(w http.ResponseWriter, r *http.Request)
	DeleteExerciseSet(w http.ResponseWriter, r *http.Request)
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

func (h *handler) CreateWorkoutSession(w http.ResponseWriter, r *http.Request) {
	var session WorkoutSession
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	session.UserID = userID

	if err := h.svc.CreateWorkoutSession(r.Context(), session); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handler) GetWorkoutSessionByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	workoutSession, err := h.svc.GetWorkoutSessionByID(r.Context(), idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workoutSession)
}

func (h *handler) ListWorkoutSessions(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	workoutSessions, err := h.svc.ListWorkoutSessions(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workoutSessions)
}

func (h *handler) UpdateWorkoutSession(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	fmt.Println("idStr: ", idStr)

	var session WorkoutSession
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	session.ID = idStr
	fmt.Println("session: ", session)

	if err := h.svc.UpdateWorkoutSession(r.Context(), session); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) DeleteWorkoutSession(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	if err := h.svc.DeleteWorkoutSession(r.Context(), idStr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) getUserIDFromContext(r *http.Request) (int32, error) {
	ctxUserID := r.Context().Value(middleware.UserIDKey)
	if ctxUserID == nil {
		return 0, errors.New("unauthorized: user id not found in context. make sure Auth middleware is active")
	}

	userIDStr, ok := ctxUserID.(string)
	if !ok {
		return 0, errors.New("internal error: user id in context is not a string")
	}

	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, fmt.Errorf("internal error: invalid user id format in context: %v", err)
	}

	return int32(userIDInt), nil
}

func (h *handler) CreateWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	var workoutExercise WorkoutExercise
	if err := json.NewDecoder(r.Body).Decode(&workoutExercise); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.svc.CreateWorkoutExercise(r.Context(), workoutExercise); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handler) GetWorkoutExerciseByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	workoutExercise, err := h.svc.GetWorkoutExerciseByID(r.Context(), idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workoutExercise)
}

func (h *handler) ListWorkoutExercisesBySession(w http.ResponseWriter, r *http.Request) {
	workoutSessionID := chi.URLParam(r, "workout_session_id")

	workoutExercises, err := h.svc.ListWorkoutExercisesBySession(r.Context(), workoutSessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workoutExercises)
}

func (h *handler) UpdateWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	var workoutExercise WorkoutExercise
	if err := json.NewDecoder(r.Body).Decode(&workoutExercise); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	workoutExercise.ID = idStr

	if err := h.svc.UpdateWorkoutExercise(r.Context(), workoutExercise); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) DeleteWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	if err := h.svc.DeleteWorkoutExercise(r.Context(), idStr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) CreateExerciseSet(w http.ResponseWriter, r *http.Request) {
	var exerciseSet ExerciseSet
	if err := json.NewDecoder(r.Body).Decode(&exerciseSet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.svc.CreateExerciseSet(r.Context(), exerciseSet); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handler) GetExerciseSetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	exerciseSet, err := h.svc.GetExerciseSetByID(r.Context(), idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exerciseSet)
}

func (h *handler) ListExerciseSets(w http.ResponseWriter, r *http.Request) {
	workoutExerciseID := chi.URLParam(r, "workout_exercise_id")

	exerciseSets, err := h.svc.ListExerciseSets(r.Context(), workoutExerciseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exerciseSets)
}

func (h *handler) UpdateExerciseSet(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	var exerciseSet ExerciseSet
	if err := json.NewDecoder(r.Body).Decode(&exerciseSet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	exerciseSet.ID = idStr

	if err := h.svc.UpdateExerciseSet(r.Context(), exerciseSet); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) DeleteExerciseSet(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	if err := h.svc.DeleteExerciseSet(r.Context(), idStr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
