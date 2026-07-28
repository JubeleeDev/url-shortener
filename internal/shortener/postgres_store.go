package shortener

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (ps *PostgresStore) Save(ctx context.Context, link Link) error {
	query := `INSERT INTO links (code, original_url) VALUES ($1, $2)`
	_, err := ps.pool.Exec(ctx, query, link.Code, link.OriginalURL)

	return err
}

func (ps *PostgresStore) Find(ctx context.Context, code string) (*Link, error) {
	query := `SELECT code, original_url FROM links WHERE code = $1`

	result := Link{}

	err := ps.pool.QueryRow(ctx, query, code).Scan(&result.Code, &result.OriginalURL)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &result, nil
}
