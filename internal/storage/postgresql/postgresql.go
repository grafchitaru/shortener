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

	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS url(
			id SERIAL PRIMARY KEY,
			user_id UUID NOT NULL,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL UNIQUE
);
		`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{conn: conn}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string, userId string) (int64, error) {
	const op = "storage.postgresql.SaveURL"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	tx, err := s.conn.Begin(ctx)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(context.Background())

	var id int64
	err = tx.QueryRow(context.Background(), `
		INSERT INTO url(url, alias, user_id) VALUES($1, $2, $3)
		ON CONFLICT (alias) DO NOTHING
		RETURNING id;
	`, urlToSave, alias, userId).Scan(&id)
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

func (s *Storage) GetUserURLs(UserId string, BaseUrl string) ([]storage.ShortURL, error) {
	const op = "storage.postgresql.GetURL"

	rows, err := s.conn.Query(context.Background(), "SELECT url, alias FROM url WHERE user_id = $1", UserId)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrURLNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	urls := make([]storage.ShortURL, 0)
	for rows.Next() {
		var url storage.ShortURL
		err := rows.Scan(&url.OriginalURL, &url.ShortURL)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		url.ShortURL = BaseUrl + url.ShortURL
		urls = append(urls, url)
	}

	return urls, nil
}

func (s *Storage) GetAlias(url string) (string, error) {
	const op = "storage.postgresql.GetAlias"

	var resAlias string
	err := s.conn.QueryRow(context.Background(), "SELECT alias FROM url WHERE url = $1", url).Scan(&resAlias)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", storage.ErrAliasNotFound
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
