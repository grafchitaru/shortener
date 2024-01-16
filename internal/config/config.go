package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

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
