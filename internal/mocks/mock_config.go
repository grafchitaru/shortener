package mocks

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/grafchitaru/shortener/internal/config"
)

func NewConfig() *config.Config {
	var cfg config.Config

	err := env.Parse(&cfg)
	if err != nil {
		fmt.Println("Can't parse  config: %w", err)
	}

	return &cfg
}
