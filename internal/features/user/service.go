package user

import (
	"context"
	"nds-go-starter/internal/core/repository"
)

type Service interface {
	ListUsers(ctx context.Context) ([]repository.MsUser, error)
}

type svc struct {
	userdb repository.Querier
}

func NewService(userdb repository.Querier) Service {
	return &svc{userdb: userdb}
}

func (s *svc) ListUsers(ctx context.Context) ([]repository.MsUser, error) {
	return s.userdb.ListUsers(ctx)
}
