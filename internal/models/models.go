package models

import "github.com/jackc/pgx/v5/pgxpool"

type Models struct {
	Users User
}

func New(db *pgxpool.Pool) *Models {
	models := &Models{
		Users: &UsersModel{
			DB: db,
		},
	}
	return models
}
