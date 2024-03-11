package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Address string
	MaxSize int
	ProdDSN string
	TestDSN string
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
		ProdDSN: os.Getenv("PROD_DSN"),
		TestDSN: os.Getenv("TEST_DSN"),
	}
	return cfg, nil
}
