package storage

type Repositories interface {
	SaveURL(urlToSave string, alias string) (int64, error)
	GetURL(alias string) (string, error)
	GetAlias(url string) (string, error)
}
