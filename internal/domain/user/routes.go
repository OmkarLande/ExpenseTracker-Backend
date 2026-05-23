package domain_user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(
	r chi.Router,
	handler *Handler,
	authMiddleware func(http.Handler) http.Handler,
) {
	r.Route("/users", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/me", handler.GetProfile)
		r.Patch("/me", handler.UpdateProfile)
		r.Post("/me/profile-picture/presign", handler.PresignProfilePicture)
		r.Post("/me/profile-picture/complete", handler.CompleteProfilePicture)
	})
}