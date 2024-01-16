package mocks

type MockStorage struct {
	SaveURLError   error
	SaveURLID      int64
	GetURLError    error
	GetURLResult   string
	GetAliasResult string
	GetAliasError  error
	PingError      error
}

func (ms *MockStorage) SaveURL(urlToSave string, alias string) (int64, error) {
	return ms.SaveURLID, ms.SaveURLError
}

func (ms *MockStorage) GetURL(alias string) (string, error) {
	return ms.GetURLResult, ms.GetURLError
}

func (ms *MockStorage) GetAlias(url string) (string, error) {
	return ms.GetAliasResult, ms.GetURLError
}

func (ms *MockStorage) Ping() error {
	return ms.PingError
}
