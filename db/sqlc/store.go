package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}
type SqlStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) *SqlStore {
	return &SqlStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
