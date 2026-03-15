package workout

import "context"

type Service interface {
	CreateExercise(ctx context.Context, exercise Exercise) error
	GetExerciseByID(ctx context.Context, id string) (Exercise, error)
	ListExercises(ctx context.Context) ([]Exercise, error)
	UpdateExercise(ctx context.Context, exercise Exercise) error

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
}

type svc struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &svc{repo: repo}
}

func (s *svc) CreateExercise(ctx context.Context, exercise Exercise) error {
	return s.repo.CreateExercise(ctx, exercise)
}

func (s *svc) GetExerciseByID(ctx context.Context, id string) (Exercise, error) {
	return s.repo.GetExerciseByID(ctx, id)
}

func (s *svc) ListExercises(ctx context.Context) ([]Exercise, error) {
	return s.repo.ListExercises(ctx)
}

func (s *svc) UpdateExercise(ctx context.Context, exercise Exercise) error {
	return s.repo.UpdateExercise(ctx, exercise)
}

func (s *svc) CreateWorkoutSession(ctx context.Context, session WorkoutSession) error {
	return s.repo.CreateWorkoutSession(ctx, session)
}

func (s *svc) GetWorkoutSessionByID(ctx context.Context, id string) (WorkoutSession, error) {
	return s.repo.GetWorkoutSessionByID(ctx, id)
}

func (s *svc) ListWorkoutSessions(ctx context.Context, userID int32) ([]WorkoutSession, error) {
	return s.repo.ListWorkoutSessions(ctx, userID)
}

func (s *svc) UpdateWorkoutSession(ctx context.Context, session WorkoutSession) error {
	return s.repo.UpdateWorkoutSession(ctx, session)
}

func (s *svc) DeleteWorkoutSession(ctx context.Context, id string) error {
	return s.repo.DeleteWorkoutSession(ctx, id)
}

func (s *svc) CreateWorkoutExercise(ctx context.Context, workoutExercise WorkoutExercise) error {
	return s.repo.CreateWorkoutExercise(ctx, workoutExercise)
}

func (s *svc) GetWorkoutExerciseByID(ctx context.Context, id string) (WorkoutExercise, error) {
	return s.repo.GetWorkoutExerciseByID(ctx, id)
}

func (s *svc) ListWorkoutExercisesBySession(ctx context.Context, workoutSessionID string) ([]WorkoutExercise, error) {
	return s.repo.ListWorkoutExercisesBySession(ctx, workoutSessionID)
}

func (s *svc) UpdateWorkoutExercise(ctx context.Context, workoutExercise WorkoutExercise) error {
	return s.repo.UpdateWorkoutExercise(ctx, workoutExercise)
}

func (s *svc) DeleteWorkoutExercise(ctx context.Context, id string) error {
	return s.repo.DeleteWorkoutExercise(ctx, id)
}

func (s *svc) CreateExerciseSet(ctx context.Context, exerciseSet ExerciseSet) error {
	return s.repo.CreateExerciseSet(ctx, exerciseSet)
}

func (s *svc) GetExerciseSetByID(ctx context.Context, id string) (ExerciseSet, error) {
	return s.repo.GetExerciseSetByID(ctx, id)
}

func (s *svc) ListExerciseSets(ctx context.Context, workoutExerciseID string) ([]ExerciseSet, error) {
	return s.repo.ListExerciseSets(ctx, workoutExerciseID)
}

func (s *svc) UpdateExerciseSet(ctx context.Context, exerciseSet ExerciseSet) error {
	return s.repo.UpdateExerciseSet(ctx, exerciseSet)
}

func (s *svc) DeleteExerciseSet(ctx context.Context, id string) error {
	return s.repo.DeleteExerciseSet(ctx, id)
}

func (s *svc) GetFullWorkoutSession(ctx context.Context, id string) (FullWorkoutSession, error) {
	return s.repo.GetFullWorkoutSession(ctx, id)
}
