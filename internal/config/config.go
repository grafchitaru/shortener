package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/grafchitaru/shortener/internal/storage"
)

type Config struct {
	HTTPServerAddress   string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseShortURL        string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
	SqliteStoragePath   string `env:"SQLITE_STORAGE_PATH" envDefault:"././internal/storage/storage.db"`
	FileDatabasePath    string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/short-url-db.json"`
	PostgresDatabaseDsn string `env:"DATABASE_DSN" envDefault:"postgres://root:root@localhost:54321/app"`
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
	flag.StringVar(&cfg.HTTPServerAddress, "a", cfg.HTTPServerAddress, "HTTP server address")
	flag.StringVar(&cfg.BaseShortURL, "b", cfg.BaseShortURL, "Base address for the resulting shortened URL")
	flag.StringVar(&cfg.FileDatabasePath, "f", cfg.FileDatabasePath, "File storage path")
	flag.StringVar(&cfg.PostgresDatabaseDsn, "d", cfg.PostgresDatabaseDsn, "PostgreSql database dsn")

	flag.Parse()

	return &cfg
}
