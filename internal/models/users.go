package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
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

var (
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrDuplicateEmail    = errors.New("duplicate email")
)

func (m *UsersModel) Insert(user *Users) (*Users, error) {
	query := `
	INSERT INTO users (id, name, username, email)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created, updated, name, username, email`

	args := []any{
		xid.New().String(),
		user.Name,
		user.Username,
		user.Email,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.Serializable,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})

	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Created,
		&user.Updated,
		&user.Name,
		&user.Username,
		&user.Email,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case strings.Contains(pgErr.Message, `duplicate key value violates unique constraint "users_username_key"`):
				return nil, ErrDuplicateUsername
			case strings.Contains(pgErr.Message, `duplicate key value violates unique constraint "users_email_key"`):
				return nil, ErrDuplicateEmail
			default:
				return nil, err
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}