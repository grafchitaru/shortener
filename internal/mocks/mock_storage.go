package mocks

import "github.com/grafchitaru/shortener/internal/storage"

type MockStorage struct {
	SaveURLError        error
	SaveURLID           int64
	GetURLError         error
	GetURLResult        string
	GetAliasResult      string
	GetAliasError       error
	PingError           error
	GetURLsResult       []storage.ShortURL
	DeleteUserURLsError error
}

func (ms *MockStorage) SaveURL(urlToSave string, alias string, userID string) (int64, error) {
	return ms.SaveURLID, ms.SaveURLError
}

func (ms *MockStorage) GetURL(alias string) (string, error) {
	return ms.GetURLResult, ms.GetURLError
}

func (ms *MockStorage) GetUserURLs(userID string, baseURL string) ([]storage.ShortURL, error) {
	return ms.GetURLsResult, ms.GetURLError
}

func (ms *MockStorage) DeleteUserURLs(userID string, deleteID []string) (string, error) {
	return "", ms.DeleteUserURLsError
}

func (ms *MockStorage) GetAlias(url string) (string, error) {
	return ms.GetAliasResult, ms.GetURLError
}

func (ms *MockStorage) Ping() error {
	return ms.PingError
}

func (ms *MockStorage) Close() {

}
