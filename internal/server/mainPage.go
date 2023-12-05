package server

import (
	"github.com/grafchitaru/shortener/internal/handlers"
	"github.com/grafchitaru/shortener/internal/storage/sqlite"
	"net/http"
)

func MainPage(res http.ResponseWriter, req *http.Request) {
	storage, err := sqlite.New("././internal/storage/storage.db")
	if err != nil {
		panic(err)
	}

	if req.Method == http.MethodPost {
		handlers.CreateLink(res, req, storage)
		return
	}

	if req.Method == http.MethodGet {
		handlers.GetLink(res, req, storage)
		return
	}

	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte("Bad Request"))
}
