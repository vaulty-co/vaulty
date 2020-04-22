package api

import (
	"net/http"

	"github.com/vaulty/proxy/storage"
)

type Server struct {
	storage storage.Storage
}

func NewServer(storage storage.Storage) *Server {
	server := &Server{
		storage: storage,
	}

	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.routes().ServeHTTP(w, r)
}

// func respondWithError(w http.ResponseWriter, code int, message string) {
//     respondWithJSON(w, code, map[string]string{"error": message})
// }

// func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
//     response, _ := json.Marshal(payload)

//     w.Header().Set("Content-Type", "application/json")
//     w.WriteHeader(code)
//     w.Write(response)
// }
