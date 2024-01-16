package handlers

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"github.com/grafchitaru/shortener/internal/app"
	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/storage"
	"io"
	"net/http"
)

func CreateLinkBatch(ctx config.HandlerContext, res http.ResponseWriter, req *http.Request) {
	var reader io.Reader
	var body []storage.BatchURL

	if req.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = req.Body
	}

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&body); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var result []storage.BatchResult

	for _, b := range body {
		alias, err := ctx.Repos.GetAlias(b.OriginalURL)
		if err != nil && !errors.Is(err, storage.ErrURLNotFound) {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		if alias != "" {
			res.WriteHeader(http.StatusConflict)
		}

		if alias == "" {
			alias = app.NewRandomString(6)
			_, err = ctx.Repos.SaveURL(b.OriginalURL, alias)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusCreated)
		}

		result = append(result, storage.BatchResult{
			CorrelationID: b.CorrelationID,
			ShortURL:      ctx.Config.BaseShortURL + "/" + alias,
		})
	}

	res.Header().Set("Content-Type", "application/json")

	json.NewEncoder(res).Encode(result)
}
