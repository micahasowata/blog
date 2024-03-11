package models

import "github.com/jackc/pgx/v5/pgxpool"

type Models struct {
}

func New(db *pgxpool.Pool) *Models {
	models := &Models{}
	return models
}
