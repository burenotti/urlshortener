package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(pool *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{pool: pool}
}

func (p PostgresStorage) SaveShortLink(ctx context.Context, linkID string, url string) error {
	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "INSERT INTO links (link_id, url) VALUES ($1, $2)", linkID, url)

	return err
}

func (p PostgresStorage) GetSource(ctx context.Context, linkID string) (string, error) {
	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, "SELECT url FROM links WHERE link_id = $1", linkID)
	var url string
	err = row.Scan(&url)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNoSourceUrl
	} else {
		return url, err
	}
}
