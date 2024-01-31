package storage

type BatchURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResult struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type Repositories interface {
	SaveURL(urlToSave string, alias string) (int64, error)
	GetURL(alias string) (string, error)
	GetAlias(url string) (string, error)
	Ping() error
}
