package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/grafchitaru/shortener/internal/storage"
	"os"
	"path/filepath"
)

type Storage struct {
	filePath string
}

func New(filePath string) (*Storage, error) {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if _, err := os.Create(filePath); err != nil {
			return nil, err
		}
	}

	return &Storage{filePath: filePath}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string, userID string) (int64, error) {
	type URLData struct {
		UUID        string `json:"uuid"`
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
		UserID      string `json:"user_id"`
	}

	uuid := uuid.New()
	data := URLData{
		UUID:        alias,
		ShortURL:    alias,
		OriginalURL: urlToSave,
		UserID:      uuid.String(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	f, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	if _, err := f.Write(append(jsonData, '\n')); err != nil {
		return 0, err
	}

	if err := f.Sync(); err != nil {
		return 0, err
	}

	return int64(len(jsonData)), nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	type URLData struct {
		UUID        string `json:"uuid"`
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	f, err := os.Open(s.filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var url URLData
		err := json.Unmarshal(scanner.Bytes(), &url)
		if err != nil {
			continue // Skip invalid lines
		}

		if url.UUID == alias {
			return url.OriginalURL, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", storage.ErrURLNotFound
}

func (s *Storage) GetAlias(url string) (string, error) {
	type URLData struct {
		UUID        string `json:"uuid"`
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	f, err := os.Open(s.filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var URLData URLData
		err := json.Unmarshal(scanner.Bytes(), &URLData)
		if err != nil {
			continue // Skip invalid lines
		}

		if URLData.OriginalURL == url {
			return URLData.UUID, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", storage.ErrAliasNotFound
}

func (s *Storage) GetUserURLs(UserID string, baseUrl string) ([]storage.ShortURL, error) {
	f, err := os.Open(s.filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var urls []storage.ShortURL
	for scanner.Scan() {
		var url storage.BatchURL
		err := json.Unmarshal(scanner.Bytes(), &url)
		if err != nil {
			continue
		}

		if url.UserID == UserID {
			urls = append(urls, storage.ShortURL{
				ShortURL:    baseUrl + url.CorrelationID,
				OriginalURL: url.OriginalURL,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func (s *Storage) Ping() error {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %w", err)
	}
	return nil
}
