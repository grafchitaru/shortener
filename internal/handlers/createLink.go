package handlers

import (
	"compress/gzip"
	"github.com/grafchitaru/shortener/internal/app"
	"github.com/grafchitaru/shortener/internal/compress"
	"github.com/grafchitaru/shortener/internal/config"
	"io"
	"net/http"
	"strings"
)

func CreateLink(ctx config.HandlerContext, res http.ResponseWriter, req *http.Request) {
	if strings.Contains(req.Header.Get("Content-Encoding"), "gzip") {
		gr, err := gzip.NewReader(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		defer gr.Close()

		req.Body = &compress.GzipReader{ReadCloser: gr}
	}
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	alias := app.NewRandomString(6)

	ctx.Repos.SaveURL(string(reqBody), alias)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(ctx.Config.BaseShortURL + "/" + alias))
}
