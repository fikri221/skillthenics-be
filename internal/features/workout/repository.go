package workout

import (
	"context"
	"database/sql"
	"nds-go-starter/internal/core/repository"
	"nds-go-starter/internal/json"

	"github.com/segmentio/ksuid"
)

func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func fromNullString(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func toNullInt32(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}

func fromNullInt32(ni sql.NullInt32) *int32 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int32
}

type Repository interface {
	CreateExercise(ctx context.Context, exercise Exercise) error
	GetExerciseByID(ctx context.Context, id string) (Exercise, error)
	ListExercises(ctx context.Context) ([]Exercise, error)
	UpdateExercise(ctx context.Context, exercise Exercise) error
	DeleteExercise(ctx context.Context, id string) error

	CreateWorkoutSession(ctx context.Context, session WorkoutSession) error
	GetWorkoutSessionByID(ctx context.Context, id string) (WorkoutSession, error)
	ListWorkoutSessions(ctx context.Context, userID int32) ([]WorkoutSession, error)
	UpdateWorkoutSession(ctx context.Context, session WorkoutSession) error
	DeleteWorkoutSession(ctx context.Context, id string) error

	CreateWorkoutExercise(ctx context.Context, workoutExercise WorkoutExercise) error
	GetWorkoutExerciseByID(ctx context.Context, id string) (WorkoutExercise, error)
	ListWorkoutExercisesBySession(ctx context.Context, workoutSessionID string) ([]WorkoutExercise, error)
	UpdateWorkoutExercise(ctx context.Context, workoutExercise WorkoutExercise) error
	DeleteWorkoutExercise(ctx context.Context, id string) error

	CreateExerciseSet(ctx context.Context, exerciseSet ExerciseSet) error
	GetExerciseSetByID(ctx context.Context, id string) (ExerciseSet, error)
	ListExerciseSets(ctx context.Context, workoutExerciseID string) ([]ExerciseSet, error)
	UpdateExerciseSet(ctx context.Context, exerciseSet ExerciseSet) error
	DeleteExerciseSet(ctx context.Context, id string) error

	GetFullWorkoutSession(ctx context.Context, id string) (FullWorkoutSession, error)
	ListFullWorkoutSessions(ctx context.Context, userID int32) ([]FullWorkoutSession, error)
}

func NewRepository(db repository.Querier) Repository {
	return &repoWrapper{db: db}
}

type repoWrapper struct {
	db repository.Querier
}

func (r *repoWrapper) CreateExercise(ctx context.Context, exercise Exercise) error {
	_, err := r.db.CreateExercise(ctx, repository.CreateExerciseParams{
		ID:          ksuid.New().String(),
		Name:        exercise.Name,
		Description: toNullString(exercise.Description),
		MuscleGroup: toNullString(exercise.MuscleGroup),
		Difficulty:  toNullString(exercise.Difficulty),
	})
	return err
}

func (r *repoWrapper) GetExerciseByID(ctx context.Context, id string) (Exercise, error) {
	exercise, err := r.db.GetExerciseByID(ctx, id)
	if err != nil {
		return Exercise{}, err
	}
	return Exercise{
		ID:          exercise.ID,
		Name:        exercise.Name,
		Description: fromNullString(exercise.Description),
		MuscleGroup: fromNullString(exercise.MuscleGroup),
		Difficulty:  fromNullString(exercise.Difficulty),
	}, nil
}

func (r *repoWrapper) ListExercises(ctx context.Context) ([]Exercise, error) {
	exercises, err := r.db.ListExercises(ctx)
	if err != nil {
		return nil, err
	}

	var result []Exercise
	for _, exercise := range exercises {
		result = append(result, mapToExercise(exercise))
	}
	return result, nil
}

func mapToExercise(e repository.Exercise) Exercise {
	return Exercise{
		ID:          e.ID,
		Name:        e.Name,
		Description: fromNullString(e.Description),
		MuscleGroup: fromNullString(e.MuscleGroup),
		Difficulty:  fromNullString(e.Difficulty),
	}
}

func (r *repoWrapper) UpdateExercise(ctx context.Context, exercise Exercise) error {
	_, err := r.db.UpdateExercise(ctx, repository.UpdateExerciseParams{
		ID:          exercise.ID,
		Name:        exercise.Name,
		Description: toNullString(exercise.Description),
		MuscleGroup: toNullString(exercise.MuscleGroup),
		Difficulty:  toNullString(exercise.Difficulty),
	})
	return err
}

func (r *repoWrapper) DeleteExercise(ctx context.Context, id string) error {
	return nil
}

func (r *repoWrapper) CreateWorkoutSession(ctx context.Context, session WorkoutSession) error {
	_, err := r.db.CreateWorkoutSession(ctx, repository.CreateWorkoutSessionParams{
		ID:              ksuid.New().String(),
		UserID:          session.UserID,
		SessionDate:     session.SessionDate,
		DurationMinutes: toNullInt32(session.Duration),
		CaloriesBurned:  toNullInt32(session.Calories),
		Notes:           toNullString(session.Notes),
	})
	return err
}

func (r *repoWrapper) GetWorkoutSessionByID(ctx context.Context, id string) (WorkoutSession, error) {
	workoutSession, err := r.db.GetWorkoutSessionByID(ctx, id)
	if err != nil {
		return WorkoutSession{}, err
	}
	return WorkoutSession{
		ID:          workoutSession.ID,
		UserID:      workoutSession.UserID,
		SessionDate: workoutSession.SessionDate,
		Duration:    fromNullInt32(workoutSession.DurationMinutes),
		Calories:    fromNullInt32(workoutSession.CaloriesBurned),
		Notes:       fromNullString(workoutSession.Notes),
	}, nil
}

func (r *repoWrapper) ListWorkoutSessions(ctx context.Context, userID int32) ([]WorkoutSession, error) {
	workoutSessions, err := r.db.ListWorkoutSessionsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	var result []WorkoutSession
	for _, workoutSession := range workoutSessions {
		result = append(result, mapToWorkoutSession(workoutSession))
	}
	return result, nil
}

func mapToWorkoutSession(workoutSession repository.WorkoutSession) WorkoutSession {
	return WorkoutSession{
		ID:          workoutSession.ID,
		UserID:      workoutSession.UserID,
		SessionDate: workoutSession.SessionDate,
		Duration:    fromNullInt32(workoutSession.DurationMinutes),
		Calories:    fromNullInt32(workoutSession.CaloriesBurned),
		Notes:       fromNullString(workoutSession.Notes),
	}
}

func (r *repoWrapper) UpdateWorkoutSession(ctx context.Context, session WorkoutSession) error {
	_, err := r.db.UpdateWorkoutSession(ctx, repository.UpdateWorkoutSessionParams{
		ID:              session.ID,
		SessionDate:     session.SessionDate,
		DurationMinutes: toNullInt32(session.Duration),
		CaloriesBurned:  toNullInt32(session.Calories),
		Notes:           toNullString(session.Notes),
	})
	return err
}

func (r *repoWrapper) DeleteWorkoutSession(ctx context.Context, id string) error {
	return nil
}

func (r *repoWrapper) CreateWorkoutExercise(ctx context.Context, workoutExercise WorkoutExercise) error {
	_, err := r.db.CreateWorkoutExercise(ctx, repository.CreateWorkoutExerciseParams{
		ID:               ksuid.New().String(),
		WorkoutSessionID: workoutExercise.WorkoutSessionID,
		ExerciseID:       workoutExercise.ExerciseID,
		Notes:            toNullString(workoutExercise.Notes),
	})
	return err
}

func (r *repoWrapper) GetWorkoutExerciseByID(ctx context.Context, id string) (WorkoutExercise, error) {
	workoutExercise, err := r.db.GetWorkoutExerciseByID(ctx, id)
	if err != nil {
		return WorkoutExercise{}, err
	}
	return WorkoutExercise{
		ID:               workoutExercise.ID,
		WorkoutSessionID: workoutExercise.WorkoutSessionID,
		ExerciseID:       workoutExercise.ExerciseID,
		Notes:            fromNullString(workoutExercise.Notes),
	}, nil
}

func (r *repoWrapper) ListWorkoutExercisesBySession(ctx context.Context, workoutSessionID string) ([]WorkoutExercise, error) {
	workoutExercises, err := r.db.ListWorkoutExercisesBySession(ctx, workoutSessionID)
	if err != nil {
		return nil, err
	}
	var result []WorkoutExercise
	for _, workoutExercise := range workoutExercises {
		result = append(result, mapToWorkoutExercise(workoutExercise))
	}
	return result, nil
}

func mapToWorkoutExercise(workoutExercise repository.ListWorkoutExercisesBySessionRow) WorkoutExercise {
	return WorkoutExercise{
		ID:               workoutExercise.ID,
		WorkoutSessionID: workoutExercise.WorkoutSessionID,
		ExerciseID:       workoutExercise.ExerciseID,
		ExerciseName:     workoutExercise.ExerciseName,
		MuscleGroup:      fromNullString(workoutExercise.MuscleGroup),
		Notes:            fromNullString(workoutExercise.Notes),
	}
}

func (r *repoWrapper) UpdateWorkoutExercise(ctx context.Context, workoutExercise WorkoutExercise) error {
	_, err := r.db.UpdateWorkoutExercise(ctx, repository.UpdateWorkoutExerciseParams{
		ID:         workoutExercise.ID,
		ExerciseID: workoutExercise.ExerciseID,
		Notes:      toNullString(workoutExercise.Notes),
	})
	return err
}

func (r *repoWrapper) DeleteWorkoutExercise(ctx context.Context, id string) error {
	return nil
}

func (r *repoWrapper) CreateExerciseSet(ctx context.Context, exerciseSet ExerciseSet) error {
	_, err := r.db.CreateExerciseSet(ctx, repository.CreateExerciseSetParams{
		ID:                ksuid.New().String(),
		WorkoutExerciseID: exerciseSet.WorkoutExerciseID,
		SetNumber:         exerciseSet.SetNumber,
		Reps:              toNullInt32(exerciseSet.Reps),
		Weight:            toNullString(exerciseSet.Weight),
		WeightUnit:        toNullString(exerciseSet.WeightUnit),
		RestSeconds:       toNullInt32(exerciseSet.RestSeconds),
	})
	return err
}

func (r *repoWrapper) GetExerciseSetByID(ctx context.Context, id string) (ExerciseSet, error) {
	exerciseSet, err := r.db.GetExerciseSetByID(ctx, id)
	if err == sql.ErrNoRows {
		return ExerciseSet{}, json.ErrNotFound
	} else if err != nil {
		return ExerciseSet{}, err
	}
	return ExerciseSet{
		ID:                exerciseSet.ID,
		WorkoutExerciseID: exerciseSet.WorkoutExerciseID,
		SetNumber:         exerciseSet.SetNumber,
		Reps:              fromNullInt32(exerciseSet.Reps),
		Weight:            fromNullString(exerciseSet.Weight),
		WeightUnit:        fromNullString(exerciseSet.WeightUnit),
		RestSeconds:       fromNullInt32(exerciseSet.RestSeconds),
	}, nil
}

func (r *repoWrapper) ListExerciseSets(ctx context.Context, workoutExerciseID string) ([]ExerciseSet, error) {
	exerciseSets, err := r.db.ListSetsByWorkoutExercise(ctx, workoutExerciseID)
	if err != nil {
		return []ExerciseSet{}, err
	}
	result := []ExerciseSet{}
	for _, exerciseSet := range exerciseSets {
		result = append(result, ExerciseSet{
			ID:                exerciseSet.ID,
			WorkoutExerciseID: exerciseSet.WorkoutExerciseID,
			SetNumber:         exerciseSet.SetNumber,
			Reps:              fromNullInt32(exerciseSet.Reps),
			Weight:            fromNullString(exerciseSet.Weight),
			WeightUnit:        fromNullString(exerciseSet.WeightUnit),
			RestSeconds:       fromNullInt32(exerciseSet.RestSeconds),
		})
	}
	return result, nil
}

func (r *repoWrapper) UpdateExerciseSet(ctx context.Context, exerciseSet ExerciseSet) error {
	_, err := r.db.UpdateExerciseSet(ctx, repository.UpdateExerciseSetParams{
		ID:          exerciseSet.ID,
		SetNumber:   exerciseSet.SetNumber,
		Reps:        toNullInt32(exerciseSet.Reps),
		Weight:      toNullString(exerciseSet.Weight),
		WeightUnit:  toNullString(exerciseSet.WeightUnit),
		RestSeconds: toNullInt32(exerciseSet.RestSeconds),
	})
	return err
}

func (r *repoWrapper) DeleteExerciseSet(ctx context.Context, id string) error {
	return nil
}

func (r *repoWrapper) GetFullWorkoutSession(ctx context.Context, id string) (FullWorkoutSession, error) {
	return FullWorkoutSession{}, nil
}

func (r *repoWrapper) ListFullWorkoutSessions(ctx context.Context, userID int32) ([]FullWorkoutSession, error) {
	return []FullWorkoutSession{}, nil
}
