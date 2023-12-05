package handlers

import (
	"github.com/grafchitaru/shortener/internal/storage/sqlite"
	"net/http"
	"strings"
)

func GetLink(res http.ResponseWriter, req *http.Request) {
	storage, err := sqlite.New("././internal/storage/storage.db")
	if err != nil {
		panic(err)
	}

	path := strings.TrimPrefix(req.URL.Path, "/")
	alias, err := storage.GetURL(path)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Bad Request"))
	}

	res.Header().Set("Location", alias)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
