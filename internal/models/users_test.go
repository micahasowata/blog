package models

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/micahasowata/blog/internal/config"
	"github.com/micahasowata/blog/internal/db"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	cfg, err := config.New()
	require.Nil(t, err)

	tdb, err := db.NewTest(cfg)

	require.Nil(t, err)
	require.NotNil(t, tdb)
	require.NotEmpty(t, tdb)

	err = db.Clean(tdb)
	require.Nil(t, err)

	return tdb
}

func TestInsert(t *testing.T) {
	tdb := setupDB(t)
	defer db.Clean(tdb)

	model := &UsersModel{
		DB: tdb,
	}

	t.Run("valid", func(t *testing.T) {
		user := &Users{
			ID:       xid.New().String(),
			Name:     "Adam",
			Username: "iamadam",
			Email:    "adam45@gmail.com",
		}

		createdUser, err := model.Insert(user)

		require.Nil(t, err)
		require.NotNil(t, user)
		require.NotEmpty(t, user)

		assert.Equal(t, user.Name, createdUser.Name)
		assert.Equal(t, user.Username, createdUser.Username)
		assert.Equal(t, user.Email, createdUser.Email)
	})

	t.Run("duplicate username", func(t *testing.T) {
		user := &Users{
			ID:       xid.NewWithTime(time.Now().Add(2 * time.Hour)).String(),
			Name:     "Adam",
			Username: "iamadam",
			Email:    "adam45@gmail.com",
		}

		createdUser, err := model.Insert(user)
		require.NotNil(t, err)
		require.Nil(t, createdUser)

		assert.EqualError(t, err, ErrDuplicateUsername.Error())
	})

	t.Run("duplicate email", func(t *testing.T) {
		user := &Users{
			ID:       xid.NewWithTime(time.Now().Add(3 * time.Hour)).String(),
			Name:     "Adam",
			Username: "iamadamthefirst",
			Email:    "adam45@gmail.com",
		}

		createdUser, err := model.Insert(user)
		require.NotNil(t, err)
		require.Nil(t, createdUser)

		assert.EqualError(t, err, ErrDuplicateEmail.Error())
	})
}

func TestVerifyEmail(t *testing.T) {
	tdb := setupDB(t)
	defer db.Clean(tdb)

	model := &UsersModel{
		DB: tdb,
	}

	user := &Users{
		ID:       xid.New().String(),
		Name:     "Adam",
		Username: "iamadam",
		Email:    "adam45@gmail.com",
	}

	createdUser, err := model.Insert(user)
	require.Nil(t, err)
	require.NotNil(t, createdUser)
	require.NotEmpty(t, createdUser)

	assert.Empty(t, createdUser.Verified)

	t.Run("valid", func(t *testing.T) {
		user, err := model.VerifyEmail(createdUser.ID)
		require.Nil(t, err)
		require.NotNil(t, user)

		require.NotEmpty(t, user.Verified)
	})

	t.Run("invalid", func(t *testing.T) {
		user, err := model.VerifyEmail(xid.New().String())
		require.NotNil(t, err)
		require.Nil(t, user)

		assert.EqualError(t, err, ErrUserNotFound.Error())
	})
}
