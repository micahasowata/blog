package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Address      string
	MaxSize      int
	ProdDSN      string
	AsynqRDC     string
	From         string
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
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

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Address:      os.Getenv("ADDR"),
		MaxSize:      size,
		ProdDSN:      os.Getenv("PROD_DSN"),
		AsynqRDC:     os.Getenv("ASYNQ_RDC"),
		From:         os.Getenv("FROM"),
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     port,
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
	}
	return cfg, nil
}
