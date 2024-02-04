package inmemory

import (
	"github.com/grafchitaru/shortener/internal/storage"
	"sync"
)

type Repositories struct {
	storage map[string]string
	mu      sync.RWMutex
}

func New() *Repositories {
	return &Repositories{
		storage: make(map[string]string),
	}
}

func (r *Repositories) SaveURL(urlToSave string, alias string, userID string) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[alias] = urlToSave
	return int64(len(r.storage)), nil
}

func (r *Repositories) GetURL(alias string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, ok := r.storage[alias]
	if !ok {
		return "", storage.ErrURLNotFound
	}
	return url, nil
}

func (r *Repositories) GetUserURLs(UserID string, baseURL string) ([]storage.ShortURL, error) {
	return nil, nil
}

func (r *Repositories) GetAlias(url string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for k, v := range r.storage {
		if v == url {
			return k, nil
		}
	}

	return "", storage.ErrAliasNotFound
}

func (r *Repositories) Ping() error {
	return nil
}
