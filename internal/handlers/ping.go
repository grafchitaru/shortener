package handlers

import (
	"github.com/grafchitaru/shortener/internal/config"
	"net/http"
)

func Ping(ctx config.HandlerContext, res http.ResponseWriter) {
	err := ctx.Repos.Ping()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}
