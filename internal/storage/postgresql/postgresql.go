package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/grafchitaru/shortener/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(connString string) (*Storage, error) {
	const op = "storage.postgresql.NewStorage"

	config, err := pgxpool.ParseConfig(connString)

	if err != nil {
		return nil, fmt.Errorf("%s: unable to parse config: %w", op, err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("%s: unable to connect: %w", op, err)
	}

	_, err = pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS url(
			id SERIAL PRIMARY KEY,
			user_id UUID NOT NULL,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL UNIQUE,
			is_deleted BOOLEAN DEFAULT FALSE
		);
	`)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{pool: pool}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string, userID string) (int64, error) {
	const op = "storage.postgresql.SaveURL"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	tx, err := s.pool.Begin(ctx)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(context.Background())

	var id int64
	err = tx.QueryRow(context.Background(), `
		INSERT INTO url(url, alias, user_id) VALUES($1, $2, $3)
		ON CONFLICT (alias) DO NOTHING
		RETURNING id;
	`, urlToSave, alias, userID).Scan(&id)
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
	var isDeleted bool
	err := s.pool.QueryRow(context.Background(), "SELECT url, is_deleted FROM url WHERE alias = $1", alias).Scan(&resURL, &isDeleted)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if isDeleted {
		return "", fmt.Errorf("isDeleted")
	}

	return resURL, nil
}

func (s *Storage) GetUserURLs(UserID string, BaseURL string) ([]storage.ShortURL, error) {
	const op = "storage.postgresql.GetUserURLs"

	rows, err := s.pool.Query(context.Background(), "SELECT url, alias FROM url WHERE user_id = $1", UserID)
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
		url.ShortURL = BaseURL + url.ShortURL
		urls = append(urls, url)
	}

	return urls, nil
}

func (s *Storage) DeleteUserURLs(userID string, deleteIDs []string) (string, error) {
	const op = "storage.postgresql.DeleteUserURLs"

	idChannel := make(chan string, len(deleteIDs))
	resultChannel := make(chan int, len(deleteIDs))

	go func() {
		defer close(idChannel)
		for _, id := range deleteIDs {
			idChannel <- id
		}
	}()

	go func() {
		for id := range idChannel {
			_, err := s.pool.Exec(context.Background(), `
                UPDATE url SET is_deleted = TRUE
                WHERE alias = $1 AND user_id = $2;
            `, id, userID)
			if err != nil {
				// TODO: logging
				continue
			}
			resultChannel <- 1
		}
	}()

	var totalDeleted int
	for range deleteIDs {
		totalDeleted += <-resultChannel
	}

	if totalDeleted == 0 {
		return "", fmt.Errorf("%s: no URLs were deleted", op)
	}

	return fmt.Sprintf("Deleted %d URLs", totalDeleted), nil
}

func (s *Storage) GetAlias(url string) (string, error) {
	const op = "storage.postgresql.GetAlias"

	var resAlias string
	err := s.pool.QueryRow(context.Background(), "SELECT alias FROM url WHERE url = $1", url).Scan(&resAlias)
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
	return s.pool.Ping(ctx)
}

func (s *Storage) Close() {
	s.pool.Close()
}
