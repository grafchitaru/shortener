package handlers

import (
	"encoding/json"
	"github.com/grafchitaru/shortener/internal/auth"
	"github.com/grafchitaru/shortener/internal/config"
	"net/http"
)

func GetUserUrls(ctx config.HandlerContext, res http.ResponseWriter, req *http.Request) {
	userId, err := auth.GetUserId(req, ctx.Config.SecretKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	baseUrl := ctx.Config.BaseShortURL + "/"
	aliases, err := ctx.Repos.GetUserURLs(userId, baseUrl)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNoContent)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(aliases)
}
