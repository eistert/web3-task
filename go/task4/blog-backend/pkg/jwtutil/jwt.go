package jwtutil

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Issue(userID uint, username string, secret []byte, ttlSeconds int64) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(ttlSeconds) * time.Second)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

// ParseToken(tokenString)：鉴权时使用；
func Parse(token string, secret []byte) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if c, ok := parsed.Claims.(*Claims); ok && parsed.Valid {
		return c, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
