package domain_auth

import (
	"net/http"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	r.Post("/auth/register", handler.Register)
	r.Post("/auth/login", handler.Login)
	r.Post("/auth/refresh", handler.RefreshToken)
	r.Post("/auth/verify-email", handler.VerifyEmail)
	r.Post("/auth/forgot-password", handler.ForgotPassword)
	r.Post("/auth/reset-password", handler.ResetPassword)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/auth/logout", handler.Logout)
	})
}
