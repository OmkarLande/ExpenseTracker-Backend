package middleware

import (
	"context"
	"net/http"
	"strings"

	"ExpenseTracker-Backend/internal/types"
	"ExpenseTracker-Backend/internal/utils"
)

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, "authorization header missing", http.StatusUnauthorized)
				return
			}

			splitToken := strings.Split(authHeader, "Bearer ")

			if len(splitToken) != 2 {
				http.Error(w, "invalid authorization format", http.StatusUnauthorized)
				return
			}

			tokenString := splitToken[1]

			claims, err := utils.ValidateAccessToken(
				tokenString,
				jwtSecret,
			)

			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				types.UserIDContextKey,
				claims.UserID,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}