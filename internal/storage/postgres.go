package storage

import (
	"context"
	"database/sql"
)

// TODO: PostgresStore through pgxpool.
// Error mapping: errors.Is(err, pgx.ErrNoRows) > ErrNotFound, pgErr.Code == "23505" > ErrCodeExist
type PostgresStore struct {
	db *sql.DB
}

func (p *PostgresStore) Save(ctx context.Context, code, longUrl string) error {
	_, err := p.db.ExecContext(ctx, `INSERT INTO urls(short_code, long_url) VALUES($1, $2) RETURNING id`, code, longUrl)

	if err != nil {
		return ErrCodeExists
	}

	return nil
}

func (p *PostgresStore) Get(ctx context.Context, code string) (string, error) {
	var longUrl string

	err := p.db.QueryRowContext(ctx, `SELECT long_url FROM urls WHERE short_code = $1`, code).Scan(&longUrl)

	if err != nil {
		return "", ErrNotFound
	}

	return longUrl, nil
}
