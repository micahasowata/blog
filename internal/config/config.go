package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	Address      string
	MaxSize      int
	ProdDSN      string
	TestDSN      string
	RDB          string
	From         string
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	Key          []byte
	IPKey        string
}

func New() (*Config, error) {
	size, err := strconv.Atoi(os.Getenv("MAX_SIZE"))
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return nil, err
	}

	key := os.Getenv("KEY")
	if len(key) != 32 {
		return nil, errors.New("token key is invalid len" + string(rune(len(key))))
	}

	cfg := &Config{
		Address:      os.Getenv("ADDR"),
		MaxSize:      size,
		ProdDSN:      os.Getenv("PROD_DSN"),
		TestDSN:      os.Getenv("TEST_DSN"),
		RDB:          os.Getenv("RDB"),
		From:         os.Getenv("FROM"),
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     port,
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		Key:          []byte(key),
		IPKey:        os.Getenv("IP_KEY"),
	}
	return cfg, nil
}
