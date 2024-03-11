package models

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/micahasowata/blog/internal/db"
	"github.com/stretchr/testify/assert"
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
	testDB := setupDB(t)
	defer db.Clean(testDB)

	model := &UsersModel{
		DB: testDB,
	}

	user := &Users{
		Name:     "Adam",
		Username: "iamadam",
		Email:    "adam45@gmail.com",
	}

	t.Run("valid", func(t *testing.T) {
		createdUser, err := model.Insert(user)

		require.Nil(t, err)
		require.NotNil(t, user)
		require.NotEmpty(t, user)

		assert.Equal(t, user.Name, createdUser.Name)
		assert.Equal(t, user.Username, createdUser.Username)
		assert.Equal(t, user.Email, createdUser.Email)
	})

	t.Run("duplicate username", func(t *testing.T) {
		createdUser, err := model.Insert(user)
		require.NotNil(t, err)
		require.Nil(t, createdUser)

		assert.EqualError(t, err, ErrDuplicateUsername.Error())
	})

	t.Run("duplicate email", func(t *testing.T) {
		user.Username = "evelyn"

		createdUser, err := model.Insert(user)
		require.NotNil(t, err)
		require.Nil(t, createdUser)

		assert.EqualError(t, err, ErrDuplicateEmail.Error())
	})
}
