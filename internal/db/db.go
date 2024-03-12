package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/micahasowata/blog/internal/config"
)

func NewProduction(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, cfg.ProdDSN)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewTest(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, cfg.TestDSN)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Clean(db *pgxpool.Pool) error {
	query := `DELETE FROM users`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
