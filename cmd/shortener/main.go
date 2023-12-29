package main

import (
	"fmt"
	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/server"
	storage2 "github.com/grafchitaru/shortener/internal/storage"
	"github.com/grafchitaru/shortener/internal/storage/file"
	"github.com/grafchitaru/shortener/internal/storage/inmemory"
	"github.com/grafchitaru/shortener/internal/storage/sqlite"
	"os"
	"path/filepath"
)

func main() {
	cfg := *config.NewConfig()

	var storage storage2.Repositories
	var err error

	if cfg.UseDatabaseFile {
		dir := filepath.Dir(cfg.FileDatabasePath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				fmt.Println("Error initialize storage: %w", err)
			}
		}

		if _, err := os.Stat(cfg.FileDatabasePath); os.IsNotExist(err) {
			if _, err := os.Create(cfg.FileDatabasePath); err != nil {
				fmt.Println("Error initialize storage: %w", err)
			}
		}
		storage, err = file.New(cfg.FileDatabasePath)
	} else if cfg.UseSqlite {
		storage, err = sqlite.New(cfg.SqliteStoragePath)
	} else {
		storage = inmemory.New()
	}
	if err != nil {
		fmt.Println("Error initialize storage: %w", err)
	}

	server.New(config.HandlerContext{Config: cfg, Repos: storage})
}
