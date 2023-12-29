package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/grafchitaru/shortener/internal/storage"
)

type Config struct {
	HTTPServerAddress string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseShortURL      string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
	UseSqlite         bool   `env:"USE_SQLITE" envDefault:"false"`
	SqliteStoragePath string `env:"SQLITE_STORAGE_PATH" envDefault:"././internal/storage/storage.db"`
	UseDatabaseFile   bool   `env:"USE_DATABASE_FILE" envDefault:"true"`
	FileDatabasePath  string `env:"FILE_STORAGE_PATH" envDefault:"././internal/storage/database.txt"`
}

type HandlerContext struct {
	Config Config
	Repos  storage.Repositories
}

func NewConfig() *Config {
	var cfg Config

	err := env.Parse(&cfg)
	if err != nil {
		fmt.Println("Can't parse  config: %w", err)
	}
	flag.StringVar(&cfg.HTTPServerAddress, "a", "127.0.0.1:8080", "HTTP server address")
	flag.StringVar(&cfg.BaseShortURL, "b", "http://127.0.0.1:8080", "Base address for the resulting shortened URL")
	flag.StringVar(&cfg.FileDatabasePath, "f", "././internal/storage/database.txt", "File storage path")

	flag.Parse()

	return &cfg
}
