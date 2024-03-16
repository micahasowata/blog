package main

import (
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccessToken(t *testing.T) {
	claims := &tokenClaims{
		ID: xid.New().String(),
	}

	app := setupApp(t, nil)

	token, err := app.newAccessToken(claims)
	require.Nil(t, err)
	require.NotEmpty(t, token)
}

func TestVerifyToken(t *testing.T) {
	app := setupApp(t, nil)

	t.Run("access token", func(t *testing.T) {
		claims := &tokenClaims{
			ID: xid.New().String(),
		}

		token, err := app.newAccessToken(claims)
		require.Nil(t, err)

		claimsFromToken, err := app.verifyJWT(token)
		require.Nil(t, err)

		assert.Equal(t, claims.ID, claimsFromToken.ID)
		assert.Equal(t, claims.StdClaims.Subject, "access")
	})

	t.Run("refresh token", func(t *testing.T) {
		claims := &tokenClaims{
			ID: xid.New().String(),
		}

		token, err := app.newRefreshToken(claims)
		require.Nil(t, err)

		claimsFromToken, err := app.verifyJWT(token)
		require.Nil(t, err)

		assert.Equal(t, claimsFromToken.StdClaims.Subject, "refresh")
	})
}
