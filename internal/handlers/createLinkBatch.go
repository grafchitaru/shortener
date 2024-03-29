package handlers

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/grafchitaru/shortener/internal/app"
	"github.com/grafchitaru/shortener/internal/auth"
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

	res.Header().Set("Content-Type", "application/json")

	for _, b := range body {
		alias, err := ctx.Repos.GetAlias(b.OriginalURL)
		if err != nil && !errors.Is(err, storage.ErrAliasNotFound) {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		userID, err := auth.GetUserID(req, ctx.Config.SecretKey)
		if err != nil {
			userID = uuid.New().String()
		}

		if alias != "" {
			res.WriteHeader(http.StatusConflict)
		}

		if alias == "" {
			alias = app.NewRandomString(6)
			_, err = ctx.Repos.SaveURL(b.OriginalURL, alias, userID)
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

	json.NewEncoder(res).Encode(result)
}
