package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type LogDB struct {
	db DBTX
}

func NewLogDB(db DBTX) *LogDB {
	return &LogDB{db: db}
}

func (l *LogDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	db, ok := l.db.(Transactor)
	if !ok {
		return nil, errors.New("underlying db does not support transactions")
	}
	return db.BeginTx(ctx, opts)
}

func (l *LogDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	res, err := l.db.ExecContext(ctx, query, args...)
	l.log(ctx, query, args, time.Since(start), err)
	return res, err
}

func (l *LogDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return l.db.PrepareContext(ctx, query)
}

func (l *LogDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := l.db.QueryContext(ctx, query, args...)
	l.log(ctx, query, args, time.Since(start), err)
	return rows, err
}

func (l *LogDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()
	row := l.db.QueryRowContext(ctx, query, args...)
	l.log(ctx, query, args, time.Since(start), nil)
	return row
}

func (l *LogDB) log(ctx context.Context, query string, args []interface{}, duration time.Duration, err error) {
	traceID := middleware.GetReqID(ctx)
	if err != nil {
		slog.Error("SQL Query Error", "traceId", traceID, "query", query, "args", args, "duration", duration, "err", err)
	} else {
		slog.Info("SQL Query", "traceId", traceID, "query", query, "args", args, "duration", duration)
	}
}
