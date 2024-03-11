package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User interface {
	Insert(*Users) (*Users, error)
}

type Users struct {
	ID       string    `json:"id"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type UsersModel struct {
	DB *pgxpool.Pool
}

func (m *UsersModel) Insert(user *Users) (*Users, error) {
	return nil, nil
}
