package domain_transaction

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(
	r chi.Router,
	handler *Handler,
	authMiddleware func(http.Handler) http.Handler,
) {
	r.Route("/transactions", func(r chi.Router) {
		r.Use(authMiddleware)
		
		r.Post("/", handler.Create)
		r.Get("/", handler.List)
		r.Get("/info", handler.GetInfo)
		r.Get("/analytics", handler.GetAnalytics)
		
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetByID)
			r.Patch("/", handler.Update)
			r.Patch("/status", handler.SoftDelete)
		})
	})
}
