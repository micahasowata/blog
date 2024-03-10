package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Address string
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Address: os.Getenv("ADDR"),
	}
	return cfg, nil
}
