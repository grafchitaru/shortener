package handlers

import (
	"compress/gzip"
	"errors"
	"github.com/google/uuid"
	"github.com/grafchitaru/shortener/internal/app"
	"github.com/grafchitaru/shortener/internal/auth"
	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/storage"
	"io"
	"net/http"
)

func CreateLink(ctx config.HandlerContext, res http.ResponseWriter, req *http.Request) {
	var reader io.Reader

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

	body, ioError := io.ReadAll(reader)
	if ioError != nil {
		http.Error(res, ioError.Error(), http.StatusBadRequest)
		return
	}

	originalURL := string(body)

	alias, err := ctx.Repos.GetAlias(originalURL)
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
		_, err = ctx.Repos.SaveURL(originalURL, alias, userID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusCreated)
	}

	res.Write([]byte(ctx.Config.BaseShortURL + "/" + alias))
}
