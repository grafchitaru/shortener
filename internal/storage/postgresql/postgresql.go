package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/grafchitaru/shortener/internal/storage"
	"github.com/jackc/pgx/v5"
	"time"
)

type Storage struct {
	conn *pgx.Conn
}

func New(connString string) (*Storage, error) {
	const op = "storage.postgresql.NewStorage"

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var exists bool
	err = conn.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1
			FROM  pg_class c
			JOIN  pg_namespace n ON n.oid = c.relnamespace
			WHERE c.relname = 'idx_alias'
			AND   n.nspname = 'public'
		)
	`).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !exists {
		_, err = conn.Exec(context.Background(), `
			CREATE TABLE IF NOT EXISTS url(
				id SERIAL PRIMARY KEY,
				alias TEXT NOT NULL UNIQUE,
				url TEXT NOT NULL);
			CREATE UNIQUE INDEX idx_alias ON url(alias);
			CREATE UNIQUE INDEX idx_url ON url(url);
		`)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &Storage{conn: conn}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.postgresql.SaveURL"

	tx, err := s.conn.Begin(context.Background())
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(context.Background())

	var id int64
	err = tx.QueryRow(context.Background(), `
		INSERT INTO url(url,alias) VALUES($1,$2)
		ON CONFLICT (alias) DO NOTHING
		RETURNING id;
	`, urlToSave, alias).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err := tx.QueryRow(context.Background(), "SELECT url FROM url WHERE alias=$1", alias).Scan(&urlToSave)
			if err != nil {
				return 0, fmt.Errorf("%s: %w", op, err)
			}
			return 0, fmt.Errorf("%s: URL %s already exists", op, urlToSave)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgresql.GetURL"

	var resURL string
	err := s.conn.QueryRow(context.Background(), "SELECT url FROM url WHERE alias = $1", alias).Scan(&resURL)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) GetAlias(url string) (string, error) {
	const op = "storage.postgresql.GetAlias"

	var resAlias string
	err := s.conn.QueryRow(context.Background(), "SELECT alias FROM url WHERE url = $1", url).Scan(&resAlias)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resAlias, nil
}

func (s *Storage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.conn.Ping(ctx)
}
