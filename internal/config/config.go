package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	HTTPServerAddress string `env:"SERVER_ADDRESS"`
	BaseShortURL      string `env:"BASE_URL"`
}

func NewConfig() *Config {
	var cfg Config

	err := env.Parse(&cfg)
	if err != nil {
		flag.StringVar(&cfg.HTTPServerAddress, "a", "localhost:8080", "HTTP server address")
		flag.StringVar(&cfg.BaseShortURL, "b", "http://localhost:8080", "Base address for the resulting shortened URL")
	}

	flag.Parse()

	return &cfg
}
