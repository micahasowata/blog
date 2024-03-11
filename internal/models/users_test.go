package models

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/micahasowata/blog/internal/db"
	"github.com/stretchr/testify/require"
)

func setupDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	db, err := db.NewTest()

	require.Nil(t, err)
	require.NotNil(t, db)
	require.NotEmpty(t, db)

	return db
}

func TestInsert(t *testing.T) {
	model := &UsersModel{
		DB: setupDB(t),
	}

	user := &Users{
		Name:     "Adam",
		Username: "iamadam",
		Email:    "adam45@gmail.com",
	}

	user, err := model.Insert(user)
	require.Nil(t, err)
	require.Nil(t, user)
}
