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
			r.Use(s.VaultCtx)
			r.Get("/", s.HandleVaultFind())
			r.Delete("/", s.HandleVaultDelete())
		})
	})
	return r
}
