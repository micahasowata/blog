package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User interface {
	Insert(*Users) (*Users, error)
	VerifyEmail(string) (*Users, error)
}

type Users struct {
	ID       string    `json:"id"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Verified bool      `json:"verified"`
}

type UsersModel struct {
	DB *pgxpool.Pool
}

var (
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrUserNotFound      = errors.New("user not found")
)

func (m *UsersModel) Insert(user *Users) (*Users, error) {
	query := `
	INSERT INTO users (id, name, username, email)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created, updated, name, username, email, verified`

	args := []any{
		user.ID,
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
		&user.Verified,
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

func (m *UsersModel) VerifyEmail(id string) (*Users, error) {
	query := `
	UPDATE users 
	SET verified = true, updated = now()
	WHERE id = $1
	RETURNING id, created, updated, name, username, email, verified`

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

	user := &Users{}

	err = tx.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Created,
		&user.Updated,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Verified,
	)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "no rows in result set"):
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}
