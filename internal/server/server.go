package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/shortener/internal/handlers"
	"github.com/grafchitaru/shortener/internal/storage/sqlite"
	"net/http"
)

func Server() {
	storage, err := sqlite.New("././internal/storage/storage.db")
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Get("/{id}", func(res http.ResponseWriter, req *http.Request) {
		handlers.GetLink(res, req, storage)
	})

	r.Post("/", func(res http.ResponseWriter, req *http.Request) {
		handlers.CreateLink(res, req, storage)
	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
