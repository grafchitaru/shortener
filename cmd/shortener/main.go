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

	if cfg.UseDatabaseFile {
		storage, err = file.New(cfg.FileDatabasePath)
	} else if cfg.UseSqlite {
		storage, err = sqlite.New(cfg.SqliteStoragePath)
	} else if cfg.UsePostgreSql {
		storage, err = postgresql.New(cfg.PostgresDatabaseDsn)
	} else {
		storage = inmemory.New()
	}
	if err != nil {
		fmt.Println("Error initialize storage: %w", err)
	}

	server.New(config.HandlerContext{Config: cfg, Repos: storage})
}
