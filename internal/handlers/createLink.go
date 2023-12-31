package handlers

import (
	"github.com/grafchitaru/shortener/internal/app"
	"github.com/grafchitaru/shortener/internal/config"
	"io"
	"net/http"
)

func CreateLink(ctx config.HandlerContext, res http.ResponseWriter, req *http.Request) {
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
