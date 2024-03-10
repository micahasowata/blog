package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Address string
	MaxSize int
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	size, err := strconv.Atoi(os.Getenv("MAX_SIZE"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Address: os.Getenv("ADDR"),
		MaxSize: size,
	}
	return cfg, nil
}
