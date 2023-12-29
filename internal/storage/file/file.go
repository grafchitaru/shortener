package file

import (
	"bufio"
	"encoding/json"
	"github.com/grafchitaru/shortener/internal/storage"
	"os"
	"path/filepath"
)

type Storage struct {
	filePath string
}

func New(filePath string) (*Storage, error) {
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return nil, err
	}
	return &Storage{filePath: filePath}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	type UrlData struct {
		Uuid        string `json:"uuid"`
		ShortUrl    string `json:"short_url"`
		OriginalUrl string `json:"original_url"`
	}

	data := UrlData{
		Uuid:        alias,
		ShortUrl:    alias,
		OriginalUrl: urlToSave,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	f, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	if _, err := f.Write(append(jsonData, '\n')); err != nil {
		return 0, err
	}

	return int64(len(jsonData)), nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	type UrlData struct {
		Uuid        string `json:"uuid"`
		ShortUrl    string `json:"short_url"`
		OriginalUrl string `json:"original_url"`
	}

	f, err := os.Open(s.filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var url UrlData
		err := json.Unmarshal(scanner.Bytes(), &url)
		if err != nil {
			continue // Skip invalid lines
		}

		if url.Uuid == alias {
			return url.OriginalUrl, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", storage.ErrURLNotFound
}

func (s *Storage) GetAlias(url string) (string, error) {
	type UrlData struct {
		Uuid        string `json:"uuid"`
		ShortUrl    string `json:"short_url"`
		OriginalUrl string `json:"original_url"`
	}

	f, err := os.Open(s.filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var urlData UrlData
		err := json.Unmarshal(scanner.Bytes(), &urlData)
		if err != nil {
			continue // Skip invalid lines
		}

		if urlData.OriginalUrl == url {
			return urlData.Uuid, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", storage.ErrAliasNotFound
}
