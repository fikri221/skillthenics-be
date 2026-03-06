package workout

import "context"

type Service interface {
	CreateExercise(ctx context.Context, exercise Exercise) error
	GetExerciseByID(ctx context.Context, id string) (Exercise, error)
	ListExercises(ctx context.Context) ([]Exercise, error)
	UpdateExercise(ctx context.Context, exercise Exercise) error
	DeleteExercise(ctx context.Context, id string) error
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

func (s *svc) DeleteExercise(ctx context.Context, id string) error {
	return s.repo.DeleteExercise(ctx, id)
}
