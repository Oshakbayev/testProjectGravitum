package v1

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func setRoutes(handler *Handler,
	mux *chi.Mux,
	authorizedJWTMdl func(http.Handler) http.Handler,
	authorizedAdminMdl func(http.Handler) http.Handler,
	recoveryMdl func(http.Handler) http.Handler,
) {
	mux.Route("/", func(r chi.Router) {
		r.Use(recoveryMdl)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", handler.Login)
		})
		r.Group(func(r chi.Router) {
			r.Use(authorizedJWTMdl)
			r.Route("/user", func(r chi.Router) {
				r.Use(authorizedAdminMdl)
				r.Get("/{id}", handler.GetUser)
				r.Post("/", handler.CreateUser)
				r.Put("/{id}", handler.UpdateUser)
			})
		})

	})
}
