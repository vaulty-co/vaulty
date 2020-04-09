package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (s *Server) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/vaults", func(r chi.Router) {
		r.Post("/", s.HandleVaultCreate())
		r.Get("/", s.HandleVaultList())
		r.Route("/{vaultID}", func(r chi.Router) {
			r.Get("/", s.HandleVaultFind())
		})
	})
	return r
}
