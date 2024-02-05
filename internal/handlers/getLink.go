package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/shortener/internal/config"
	"net/http"
	"strings"
)

func GetLink(ctx config.HandlerContext, res http.ResponseWriter, req *http.Request) {
	path := chi.URLParam(req, "id")
	alias, err := ctx.Repos.GetURL(path)
	if err != nil {
		if strings.Contains(err.Error(), "Url Is Deleted") {
			http.Error(res, "The requested resource has been deleted.", http.StatusGone)
		} else {
			http.Error(res, err.Error(), http.StatusNotFound)
		}
		return
	}

	res.Header().Set("Location", alias)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
