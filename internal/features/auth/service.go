package auth

import (
	"context"
	"nds-go-starter/internal/core/auth"
	"strconv"
	"time"

	"github.com/segmentio/ksuid"
)

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type Service interface {
	Login(ctx context.Context, email, password string) (LoginResponse, error)
	Refresh(ctx context.Context, refreshToken string) (string, error)
}

type service struct {
	repo       Repository
	jwtManager *auth.JWTManager
}

func NewService(repo Repository, jwtManager *auth.JWTManager) Service {
	return &service{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (s *service) Login(ctx context.Context, email, password string) (LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return LoginResponse{}, err
	}

	// Security: Use bcrypt to compare the hashed password from database with the provided password
	if !auth.CheckPasswordHash(password, user.Password) {
		return LoginResponse{}, auth.ErrInvalidToken
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(strconv.Itoa(int(user.ID)))
	if err != nil {
		return LoginResponse{}, err
	}

	refreshToken := ksuid.New().String()
	session := Session{
		ID:           ksuid.New().String(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (string, error) {
	session, err := s.repo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", auth.ErrInvalidToken
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.repo.DeleteSession(ctx, session.ID)
		return "", auth.ErrExpiredToken
	}

	return s.jwtManager.GenerateAccessToken(strconv.Itoa(int(session.UserID)))
}
