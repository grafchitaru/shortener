package main

import (
	"github.com/grafchitaru/shortener/internal/config"
	"github.com/grafchitaru/shortener/internal/server"
)

func main() {
	cfg := *config.NewConfig()
	server.New(cfg)
}
