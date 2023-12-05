package handlers

import (
	"github.com/grafchitaru/shortener/internal/app"
	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/storage"
	"io"
	"net/http"
)

func CreateLink(res http.ResponseWriter, req *http.Request, storage storage.Repositories, cfg *config.Config) {
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(err.Error()))
		return
	}
	alias := app.NewRandomString(6)

	storage.SaveURL(string(reqBody), alias)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(cfg.BaseShortURL + alias))
}
