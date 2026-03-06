package repository

import (
	"context"
	"database/sql"
)

type Transactor interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
