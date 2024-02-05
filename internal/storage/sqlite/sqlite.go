package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/grafchitaru/shortener/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/mattn/go-sqlite3"
	"runtime"
	"sync"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
    CREATE TABLE IF NOT EXISTS url(
        id INTEGER PRIMARY KEY,
        user_id TEXT NOT NULL,
        alias TEXT NOT NULL UNIQUE,
        url TEXT NOT NULL,
        is_deleted BOOLEAN DEFAULT FALSE);
    CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string, userID string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias, user_id) values(?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) GetUserURLs(UserID string, baseURL string) ([]storage.ShortURL, error) {
	const op = "storage.sqlite.GetURL"

	rows, err := s.db.Query("SELECT url, alias FROM url WHERE user_id = $1", UserID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrURLNotFound
	}
	defer rows.Close()

	for rows.Next() {
		// process rows here
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrURLNotFound
	}

	urls := make([]storage.ShortURL, 0)
	for rows.Next() {
		var url storage.ShortURL
		err := rows.Scan(&url.OriginalURL, &url.ShortURL)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		url.ShortURL = baseURL + url.ShortURL
		urls = append(urls, url)
	}

	return urls, nil
}

func (s *Storage) DeleteUserURLs(userID string, deleteIDs []string) (string, error) {
	const op = "storage.postgresql.DeleteUserURLs"

	idChannel := make(chan string, len(deleteIDs))
	resultChannel := make(chan bool, len(deleteIDs))

	go func() {
		for _, id := range deleteIDs {
			idChannel <- id
		}
		close(idChannel)
	}()

	numWorkers := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for id := range idChannel {
				_, err := s.db.ExecContext(context.Background(), `
                    UPDATE url SET is_deleted = TRUE
                    WHERE alias = $1 AND user_id = $2;
                `, id, userID)
				if err != nil {
					//TODO logging
					continue
				}
				resultChannel <- true
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	totalDeleted := 0
	for range resultChannel {
		totalDeleted++
	}

	if totalDeleted == 0 {
		return "", fmt.Errorf("%s: no URLs were deleted", op)
	}

	return fmt.Sprintf("Deleted %d URLs", totalDeleted), nil
}

func (s *Storage) GetAlias(url string) (string, error) {
	const op = "storage.sqlite.GetAlias"

	stmt, err := s.db.Prepare("SELECT alias FROM url WHERE url = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resAlias string

	err = stmt.QueryRow(url).Scan(&resAlias)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resAlias, nil
}

func (s *Storage) Ping() error {
	return s.db.Ping()
}
