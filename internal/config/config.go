package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	HTTPServerAddress string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseShortURL      string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
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
