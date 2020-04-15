package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/vaulty/proxy/api/request"
	"github.com/vaulty/proxy/model"
)

type vaultInput struct {
	Upstream string `json:"upstream"`
}

func (s *Server) HandleVaultCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := &model.Vault{}
		err := json.NewDecoder(r.Body).Decode(in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		vault := &model.Vault{
			Upstream: in.Upstream,
		}

		err = s.storage.CreateVault(vault)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(vault)
	}
}

func (s *Server) HandleVaultList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vaults, err := s.storage.ListVaults()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(vaults)
	}
}

func (s *Server) VaultCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vaultID := chi.URLParam(r, "vaultID")
		vault, err := s.storage.FindVault(vaultID)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		ctx := request.WithVault(r.Context(), vault)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) HandleVaultFind() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vault := request.VaultFrom(r.Context())

		json.NewEncoder(w).Encode(vault)
	}
}

func (s *Server) HandleVaultDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vault := request.VaultFrom(r.Context())
		err := s.storage.DeleteVault(vault.ID)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) HandleVaultUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vault := request.VaultFrom(r.Context())

		in := &vaultInput{}
		err := json.NewDecoder(r.Body).Decode(in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		vault.Upstream = in.Upstream

		err = s.storage.UpdateVault(vault)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(vault)
	}
}
