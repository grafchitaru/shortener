package main

import (
	"fmt"
	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/server"
	storage2 "github.com/grafchitaru/shortener/internal/storage"
	"github.com/grafchitaru/shortener/internal/storage/file"
	"github.com/grafchitaru/shortener/internal/storage/inmemory"
	"github.com/grafchitaru/shortener/internal/storage/postgresql"
	"github.com/grafchitaru/shortener/internal/storage/sqlite"
)

func main() {
	cfg := *config.NewConfig()

	var storage storage2.Repositories
	var err error

	if cfg.PostgresDatabaseDsn != "" {
		storage, err = postgresql.New(cfg.PostgresDatabaseDsn)
	} else if cfg.FileDatabasePath != "" {
		storage, err = file.New(cfg.FileDatabasePath)
	} else if cfg.SqliteStoragePath != "" {
		storage = inmemory.New()
	} else {
		storage, err = sqlite.New(cfg.SqliteStoragePath)
	}
	if err != nil {
		fmt.Println("Error initialize storage: %w", err)
	}

	defer storage.Close()

	server.New(config.HandlerContext{Config: cfg, Repos: storage})
}
