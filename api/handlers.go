package api

import (
	"encoding/json"
	"net/http"

	"github.com/vaulty/proxy/model"
)

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
