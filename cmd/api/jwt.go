package main

import (
	"time"

	"github.com/kataras/jwt"
	"github.com/rs/xid"
)

type tokenClaims struct {
	ID        string
	StdClaims *jwt.Claims
}

func (app *application) newAccessToken(claims *tokenClaims) (string, error) {
	if claims.StdClaims == nil {
		claims.StdClaims = &jwt.Claims{
			ID:        xid.New().String(),
			OriginID:  xid.New().String()[0:7] + xid.New().String(),
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			Expiry:    time.Now().Add(3 * time.Hour).Unix(),
			Issuer:    "blog-be",
			Subject:   "access",
			Audience:  jwt.Audience{"blog-ui"},
		}
	}

	token, err := jwt.Sign(jwt.HS256, app.config.Key, claims, jwt.MaxAge(4*time.Hour))
	if err != nil {
		return "", err
	}

	return string(token), nil
}

func (app *application) newRefreshToken(claims *tokenClaims) (string, error) {
	if claims.StdClaims == nil {
		claims.StdClaims = &jwt.Claims{
			ID:        xid.New().String(),
			OriginID:  xid.New().String()[0:7] + xid.New().String(),
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			Expiry:    time.Now().Add(3 * time.Hour).Unix(),
			Issuer:    "blog-be",
			Subject:   "refresh",
			Audience:  jwt.Audience{"blog-ui"},
		}
	}

	token, err := jwt.Sign(jwt.HS256, app.config.Key, claims, jwt.MaxAge(24*2*time.Hour))
	if err != nil {
		return "", err
	}

	return string(token), nil
}

type tokenPair struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

func (app *application) newTokenPair(accessToken string, refreshToken string) *tokenPair {
	pair := &tokenPair{
		Access:  accessToken,
		Refresh: refreshToken,
	}

	return pair
}

func (app *application) verifyJWT(token string) (*tokenClaims, error) {
	verifiedToken, err := jwt.Verify(jwt.HS256, app.config.Key, []byte(token))
	if err != nil {
		return nil, err
	}

	claims := &tokenClaims{}
	err = verifiedToken.Claims(&claims)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
