package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/grafchitaru/shortener/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/mattn/go-sqlite3"
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
        url TEXT NOT NULL);
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

func (s *Storage) GetUserURLs(UserID string, baseUrl string) ([]storage.ShortURL, error) {
	const op = "storage.sqlite.GetURL"

	rows, err := s.db.Query("SELECT url, alias FROM url WHERE user_id = $1", UserID)
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
		url.ShortURL = baseUrl + url.ShortURL
		urls = append(urls, url)
	}

	return urls, nil
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
