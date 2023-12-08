package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/handlers"
	storage2 "github.com/grafchitaru/shortener/internal/storage"
	"github.com/grafchitaru/shortener/internal/storage/inmemory"
	"github.com/grafchitaru/shortener/internal/storage/sqlite"
	"net/http"
)

func New(cfg config.Config) {
	var storage storage2.Repositories
	var err error

	if cfg.UseSqlite == true {
		storage, err = sqlite.New(cfg.SqliteStoragePath)
	} else {
		storage = inmemory.New()
	}
	if err != nil {
		fmt.Println("Error initialize storage: %w", err)
	}

	r := chi.NewRouter()
	r.Get("/{id}", func(res http.ResponseWriter, req *http.Request) {
		handlers.GetLink(config.HandlerContext{Config: cfg, Repos: storage}, res, req)
	})

	r.Post("/", func(res http.ResponseWriter, req *http.Request) {
		handlers.CreateLink(config.HandlerContext{Config: cfg, Repos: storage}, res, req)
	})

	err = http.ListenAndServe(cfg.HTTPServerAddress, r)
	if err != nil {
		fmt.Println("Error server: %w", err)
	}
}
