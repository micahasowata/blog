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
		user, err := model.VerifyEmail(createdUser.Email)
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

func TestGetByEmail(t *testing.T) {
	tdb := setupDB(t)
	defer db.Clean(tdb)

	model := &UsersModel{
		DB: tdb,
	}

	user := &Users{
		ID:       xid.New().String(),
		Name:     "Adam",
		Username: "iamadam47",
		Email:    "adam459@gmail.com",
	}

	createdUser, err := model.Insert(user)
	require.Nil(t, err)
	require.NotNil(t, createdUser)
	require.NotEmpty(t, createdUser)

	t.Run("valid", func(t *testing.T) {
		userFromDB, err := model.GetByEmail(createdUser.Email)
		require.Nil(t, err)
		require.NotNil(t, userFromDB)
		require.NotEmpty(t, userFromDB)

		assert.Equal(t, createdUser, userFromDB)
	})

	t.Run("invalid", func(t *testing.T) {
		userFromDB, err := model.GetByEmail("saywarawara@gmail.com")
		require.NotNil(t, err)
		require.Nil(t, userFromDB)

		assert.EqualError(t, err, ErrUserNotFound.Error())
	})
}

func TestGetByID(t *testing.T) {
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

	t.Run("valid", func(t *testing.T) {
		userFromDB, err := model.GetByID(createdUser.ID)
		require.Nil(t, err)
		require.NotNil(t, userFromDB)
		require.NotEmpty(t, userFromDB)

		assert.Equal(t, createdUser, userFromDB)
	})

	t.Run("invalid", func(t *testing.T) {
		userFromDB, err := model.GetByID(xid.New().String())
		require.NotNil(t, err)
		require.Nil(t, userFromDB)

		assert.EqualError(t, err, ErrUserNotFound.Error())
	})
}

func TestUpdateUser(t *testing.T) {
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

	user = &Users{
		ID:       xid.New().String(),
		Name:     "Jimmy",
		Username: "iamjim",
		Email:    "olaoluwa567@gmail.com",
	}

	secondUser, err := model.Insert(user)
	require.Nil(t, err)

	tests := []struct {
		name           string
		user           *Users
		shouldCauseErr bool
		err            error
	}{
		{
			name: "valid",
			user: &Users{
				ID:       createdUser.ID,
				Name:     "Victor",
				Username: createdUser.Username,
				Email:    createdUser.Email,
			},
			shouldCauseErr: false,
		},
		{
			name: "missing user",
			user: &Users{
				ID:       xid.New().String(),
				Name:     createdUser.Name,
				Username: createdUser.Username,
				Email:    "victor35@gmail.com",
			},
			shouldCauseErr: true,
			err:            ErrUserNotFound,
		},
		{
			name: "duplicate username",
			user: &Users{
				ID:       secondUser.ID,
				Name:     createdUser.Name,
				Username: createdUser.Username,
				Email:    "victor35@gmail.com",
			},
			shouldCauseErr: true,
			err:            ErrDuplicateUsername,
		},
		{
			name: "duplicate email",
			user: &Users{
				ID:       secondUser.ID,
				Name:     createdUser.Name,
				Username: "victortheGreat",
				Email:    createdUser.Email,
			},
			shouldCauseErr: true,
			err:            ErrDuplicateEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := model.Update(tt.user)
			if tt.shouldCauseErr {
				require.NotNil(t, err)
				assert.EqualError(t, err, tt.err.Error())
				assert.Nil(t, u)
			} else {
				require.Nil(t, err)
				require.NotNil(t, u)
			}
		})
	}
}
