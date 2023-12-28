package handlers

import (
	"compress/gzip"
	"github.com/grafchitaru/shortener/internal/app"
	"github.com/grafchitaru/shortener/internal/config"
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

	alias := app.NewRandomString(6)

	ctx.Repos.SaveURL(string(body), alias)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(ctx.Config.BaseShortURL + "/" + alias))
}
