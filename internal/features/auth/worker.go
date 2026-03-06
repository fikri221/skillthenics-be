package auth

import (
	"context"
	"log/slog"
	"nds-go-starter/internal/core/worker"
	"time"
)

type sessionCleanupWorker struct {
	repo     Repository
	interval time.Duration
}

func NewSessionCleanupWorker(repo Repository, interval time.Duration) worker.Worker {
	return &sessionCleanupWorker{
		repo:     repo,
		interval: interval,
	}
}

func (w *sessionCleanupWorker) Start(ctx context.Context) {
	slog.Info("Session Cleanup Worker started", "interval", w.interval)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Session Cleanup Worker stopping...")
			return
		case <-ticker.C:
			w.cleanup(ctx)
		}
	}
}

func (w *sessionCleanupWorker) cleanup(ctx context.Context) {
	slog.Info("Session Cleanup Worker: cleaning up expired sessions...")
	err := w.repo.CleanupExpiredSessions(ctx)
	if err != nil {
		slog.Error("Session Cleanup Worker: failed to cleanup sessions", "error", err)
	}
}
