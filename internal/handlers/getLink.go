package handlers

import (
	"github.com/grafchitaru/shortener/internal/storage"
	"net/http"
	"strings"
)

func GetLink(res http.ResponseWriter, req *http.Request, storage storage.Repositories) {
	path := strings.TrimPrefix(req.URL.Path, "/")
	alias, err := storage.GetURL(path)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Bad Request"))
	}

	res.Header().Set("Location", alias)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
