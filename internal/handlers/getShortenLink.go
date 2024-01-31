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

type Link struct {
	URL string `json:"url"`
}

type Result struct {
	Result string `json:"result"`
}

func GetShorten(ctx config.HandlerContext, res http.ResponseWriter, req *http.Request) {
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

	var link Link

	if err := json.Unmarshal(body, &link); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	url := link.URL

	res.Header().Set("Content-Type", "application/json")

	status := http.StatusOK
	alias, err := ctx.Repos.GetAlias(url)
	if err != nil && !errors.Is(err, storage.ErrAliasNotFound) {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if alias != "" {
		res.WriteHeader(http.StatusConflict)
	}

	if alias == "" {
		alias = app.NewRandomString(6)
		ctx.Repos.SaveURL(url, alias)
		status = http.StatusCreated
	}

	result := Result{
		Result: ctx.Config.BaseShortURL + "/" + alias,
	}
	resp, err := json.Marshal(result)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(status)
	res.Write([]byte(resp))
}
