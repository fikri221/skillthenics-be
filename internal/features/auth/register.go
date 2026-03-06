package auth

import (
	"nds-go-starter/internal/core/auth"
	"nds-go-starter/internal/core/repository"

	"github.com/go-chi/chi/v5"
)

func Register(r chi.Router, db repository.DBTX, jwtManager *auth.JWTManager) {
	querier := repository.New(db)
	repo := NewRepository(querier)
	service := NewService(repo, jwtManager)
	handler := NewHandler(service)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Post("/refresh", handler.Refresh)
	})
}
