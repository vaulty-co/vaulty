package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (s *Server) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/vaults", s.HandleVaultList())
	r.Post("/vaults", s.HandleVaultCreate())
	return r
}
