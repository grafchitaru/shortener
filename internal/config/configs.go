package config

import (
	"github.com/grafchitaru/shortener/internal/storage"
)

type Config struct {
	HTTPServerAddress   string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseShortURL        string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
	SqliteStoragePath   string `env:"SQLITE_STORAGE_PATH" envDefault:"././internal/storage/storage.db"`
	FileDatabasePath    string `env:"FILE_STORAGE_PATH"`
	PostgresDatabaseDsn string `env:"DATABASE_DSN"` // envDefault:"postgres://root:root@localhost:54321/app"
	SecretKey           string `env:"SECRET_KEY" envDefault:"your_secret_key"`
}

type HandlerContext struct {
	Config Config
	Repos  storage.Repositories
}

type Configs interface {
	NewConfig() *Config
}
