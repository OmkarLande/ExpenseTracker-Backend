package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateTokens(userID int64, secret string) (string, string, time.Time, error) {
	// Access token
	accessExpirationTime := time.Now().Add(15 * time.Minute)
	accessClaims := &TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Refresh token
	refreshExpirationTime := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := &TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessTokenString, refreshTokenString, accessExpirationTime, nil
}

func ValidateAccessToken(tokenString string, secret string) (*TokenClaims, error) {
	claims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
