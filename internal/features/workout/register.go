package workout

import (
	"nds-go-starter/internal/core/repository"

	"github.com/go-chi/chi/v5"
)

func Register(r chi.Router, db repository.DBTX) {
	querier := repository.New(db)
	repo := NewRepository(querier)
	service := NewService(repo)
	handler := NewHandler(service)

	r.Route("/workout", func(r chi.Router) {
		r.Get("/", handler.ListExercises)
		r.Post("/", handler.CreateExercise)
		r.Get("/{id}", handler.GetExerciseByID)
		r.Put("/{id}", handler.UpdateExercise)
		r.Delete("/{id}", handler.DeleteExercise)
	})

	r.Route("/workout-sessions", func(r chi.Router) {
		r.Get("/", handler.ListWorkoutSessions)
		r.Post("/", handler.CreateWorkoutSession)
		r.Get("/{id}", handler.GetWorkoutSessionByID)
		r.Put("/{id}", handler.UpdateWorkoutSession)
		r.Delete("/{id}", handler.DeleteWorkoutSession)
	})

	r.Route("/workout-exercises", func(r chi.Router) {
		r.Get("/session/{workout_session_id}", handler.ListWorkoutExercisesBySession)
		r.Post("/", handler.CreateWorkoutExercise)
		r.Get("/{id}", handler.GetWorkoutExerciseByID)
		r.Put("/{id}", handler.UpdateWorkoutExercise)
		r.Delete("/{id}", handler.DeleteWorkoutExercise)
	})

	r.Route("/exercise-sets", func(r chi.Router) {
		r.Get("/workout-exercise/{workout_exercise_id}", handler.ListExerciseSets)
		r.Post("/", handler.CreateExerciseSet)
		r.Get("/{id}", handler.GetExerciseSetByID)
		r.Put("/{id}", handler.UpdateExerciseSet)
		r.Delete("/{id}", handler.DeleteExerciseSet)
	})
}
