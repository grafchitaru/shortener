package storage

type BatchURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
	UserID        string `json:"user_id"`
}

type BatchResult struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type ShortURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Repositories interface {
	SaveURL(urlToSave string, alias string, userID string) (int64, error)
	GetURL(alias string) (string, error)
	GetUserURLs(userID string, baseURL string) ([]ShortURL, error)
	DeleteUserURLs(userID string, DeleteID []string) (string, error)
	GetAlias(url string) (string, error)
	Ping() error
}
