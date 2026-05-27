package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port   string // Default 8888
	DBPath string // Default :memory:
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	cfg := Config{}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	cfg.Port = port

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = ":memory:"
	}
	cfg.DBPath = dbPath

	return &cfg, nil
}
