package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/shortener/internal/config"
	"net/http"
)

func GetLink(ctx config.HandlerContext, res http.ResponseWriter, req *http.Request) {
	path := chi.URLParam(req, "id")
	if path == "" {
		http.Error(res, "Error read id param", http.StatusInternalServerError)
		return
	}

	alias, err := ctx.Repos.GetURL(path)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	res.WriteHeader(http.StatusTemporaryRedirect)
	res.Header().Set("Location", alias)
}
