package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PgStatStore struct {
	pool *pgxpool.Pool
}

func NewPgStatStore(pool *pgxpool.Pool) *PgStatStore {
	return &PgStatStore{
		pool: pool,
	}
}

func (p PgStatStore) RegisterRedirect(ctx context.Context, linkID, sessionID string, at time.Time) error {
	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	q := "INSERT INTO redirects (session_id, link_id, time) VALUES ($1, $2, $3)"
	_, err = conn.Exec(ctx, q, linkID, sessionID, at.UTC())
	if err != nil {
		return err
	} else {
		return nil
	}
}
