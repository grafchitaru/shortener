package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/shortener/internal/config"
	"net/http"
)

func GetLink(ctx config.HandlerContext, res http.ResponseWriter, req *http.Request) {
	path := chi.URLParam(req, "id")
	alias, err := ctx.Repos.GetURL(path)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	res.Header().Set("Location", alias)
	res.WriteHeader(http.StatusTemporaryRedirect)

}
