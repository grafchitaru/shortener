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
			url TEXT NOT NULL UNIQUE,
			is_deleted BOOLEAN DEFAULT FALSE
);
		`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{conn: conn}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string, userID string) (int64, error) {
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
	err := s.conn.QueryRow(context.Background(), "SELECT url, is_deleted FROM url WHERE alias = $1", alias).Scan(&resURL, &isDeleted)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if isDeleted == true {
		return "", fmt.Errorf("%w", "Url Is Deleted")
	}

	return resURL, nil
}

func (s *Storage) GetUserURLs(UserID string, BaseURL string) ([]storage.ShortURL, error) {
	const op = "storage.postgresql.GetUserURLs"

	rows, err := s.conn.Query(context.Background(), "SELECT url, alias FROM url WHERE user_id = $1", UserID)
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
	resultChannel := make(chan bool, len(deleteIDs))
	done := make(chan bool)

	go func() {
		for _, id := range deleteIDs {
			idChannel <- id
		}
		close(idChannel)
	}()

	go func() {
		for id := range idChannel {
			_, err := s.conn.Exec(context.Background(), `
                UPDATE url SET is_deleted = TRUE
                WHERE alias = $1 AND user_id = $2;
            `, id, userID)
			if err != nil {
				// TODO: logging
				continue
			}
			resultChannel <- true
		}
		done <- true
	}()

	totalDeleted := 0
	for i := 0; i < len(deleteIDs); i++ {
		select {
		case <-done:
			break
		case <-resultChannel:
			totalDeleted++
		}
	}

	if totalDeleted == 0 {
		return "", fmt.Errorf("%s: no URLs were deleted", op)
	}

	return fmt.Sprintf("Deleted %d URLs", totalDeleted), nil
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
