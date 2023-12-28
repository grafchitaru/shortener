package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/grafchitaru/shortener/internal/app"
	"github.com/grafchitaru/shortener/internal/compress"
	"github.com/grafchitaru/shortener/internal/config"
	"net/http"
)

type Link struct {
	URL string `json:"url"`
}

type Result struct {
	Result string `json:"result"`
}

func GetShorten(ctx config.HandlerContext, res http.ResponseWriter, req *http.Request) {
	decompressReq, err := compress.GzipDecompress(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var link Link
	var buf bytes.Buffer
	_, err = buf.ReadFrom(decompressReq.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &link); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	url := link.URL

	status := http.StatusOK
	alias, err := ctx.Repos.GetAlias(url)
	if err != nil {
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

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	res.Write([]byte(resp))
}
