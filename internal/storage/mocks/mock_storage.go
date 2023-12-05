package mocks

type MockStorage struct {
	SaveURLError error
	SaveURLID    int64
	GetURLError  error
	GetURLResult string
}

func (ms *MockStorage) SaveURL(urlToSave string, alias string) (int64, error) {
	return ms.SaveURLID, ms.SaveURLError
}

func (ms *MockStorage) GetURL(alias string) (string, error) {
	return ms.GetURLResult, ms.GetURLError
}
