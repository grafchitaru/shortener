package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/shortener/internal/storage"
	"net/http"
)

func GetLink(res http.ResponseWriter, req *http.Request, storage storage.Repositories) {
	path := chi.URLParam(req, "id")
	alias, err := storage.GetURL(path)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Bad Request"))
	}

	res.Header().Set("Location", alias)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
