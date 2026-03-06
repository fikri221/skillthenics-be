package auth

import (
	"context"
	"nds-go-starter/internal/core/repository"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (User, error)
	CreateSession(ctx context.Context, session Session) error
	GetSession(ctx context.Context, id string) (Session, error)
	GetSessionByRefreshToken(ctx context.Context, token string) (Session, error)
	DeleteSession(ctx context.Context, id string) error
	CleanupExpiredSessions(ctx context.Context) error
}

type repoWrapper struct {
	db repository.Querier
}

func NewRepository(db repository.Querier) Repository {
	return &repoWrapper{db: db}
}

func (r *repoWrapper) GetUserByEmail(ctx context.Context, email string) (User, error) {
	u, err := r.db.GetUserByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *repoWrapper) CreateSession(ctx context.Context, s Session) error {
	_, err := r.db.CreateSession(ctx, repository.CreateSessionParams{
		ID:           s.ID,
		UserID:       s.UserID,
		RefreshToken: s.RefreshToken,
		ExpiresAt:    s.ExpiresAt,
	})
	return err
}

func (r *repoWrapper) GetSession(ctx context.Context, id string) (Session, error) {
	s, err := r.db.GetSession(ctx, id)
	if err != nil {
		return Session{}, err
	}
	return Session{
		ID:           s.ID,
		UserID:       s.UserID,
		RefreshToken: s.RefreshToken,
		ExpiresAt:    s.ExpiresAt,
	}, nil
}

func (r *repoWrapper) GetSessionByRefreshToken(ctx context.Context, token string) (Session, error) {
	s, err := r.db.GetSessionByRefreshToken(ctx, token)
	if err != nil {
		return Session{}, err
	}
	return Session{
		ID:           s.ID,
		UserID:       s.UserID,
		RefreshToken: s.RefreshToken,
		ExpiresAt:    s.ExpiresAt,
	}, nil
}

func (r *repoWrapper) DeleteSession(ctx context.Context, id string) error {
	return r.db.DeleteSession(ctx, id)
}

func (r *repoWrapper) CleanupExpiredSessions(ctx context.Context) error {
	return r.db.CleanupExpiredSessions(ctx)
}
