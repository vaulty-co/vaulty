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
