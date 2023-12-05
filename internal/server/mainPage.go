package server

import (
	"github.com/grafchitaru/shortener/internal/handlers"
	"net/http"
)

func MainPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		handlers.CreateLink(res, req)
		return
	}

	if req.Method == http.MethodGet {
		handlers.GetLink(res, req)
		return
	}

	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte("Bad Request"))
}
