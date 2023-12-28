package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/shortener/internal/compress"
	"github.com/grafchitaru/shortener/internal/config"
	"net/http"
)

func GetLink(ctx config.HandlerContext, res http.ResponseWriter, req *http.Request) {
	decompressReq, err := compress.GzipDecompress(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	path := chi.URLParam(decompressReq, "id")
	alias, err := ctx.Repos.GetURL(path)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	res.Header().Set("Location", alias)
	res.WriteHeader(http.StatusTemporaryRedirect)

}
