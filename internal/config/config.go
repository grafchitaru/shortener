package config

import "flag"

type Config struct {
	HttpServerAddress string
	BaseShortURL      string
}

func NewConfig() *Config {
	var cfg Config

	flag.StringVar(&cfg.HttpServerAddress, "a", "localhost:8080", "HTTP server address")
	flag.StringVar(&cfg.BaseShortURL, "b", "http://localhost:8080/", "Base address for the resulting shortened URL")

	flag.Parse()

	return &cfg
}
