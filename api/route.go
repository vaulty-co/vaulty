package api

import (
	"encoding/json"
	"net/http"

	"github.com/vaulty/proxy/api/request"
	"github.com/vaulty/proxy/model"
)

func (s *Server) HandleRouteCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vault := request.VaultFrom(r.Context())
		route := &model.Route{}

		err := json.NewDecoder(r.Body).Decode(route)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		route.VaultID = vault.ID

		err = s.storage.CreateRoute(route)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(route)
	}
}
