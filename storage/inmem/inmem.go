package inmem

import (
	"log"
	"net/url"

	"github.com/vaulty/vaulty/model"
	"github.com/vaulty/vaulty/storage"
)

type inmemStorage struct {
	vaults         map[string]*model.Vault
	routes         map[string]*model.Route
	vaultRoutesIDs map[vaultRouteKey]*model.Route
}

// structs to store vaults and routes
type vaultRouteKey struct {
	vaultID, routeID string
}

func NewStorage() storage.Storage {
	return &inmemStorage{
		vaults:         map[string]*model.Vault{},
		routes:         map[string]*model.Route{},
		vaultRoutesIDs: map[vaultRouteKey]*model.Route{},
	}
}

func (s *inmemStorage) Reset() {
	s.routes = map[string]*model.Route{}
	s.vaultRoutesIDs = map[vaultRouteKey]*model.Route{}
	s.vaults = map[string]*model.Vault{}
}

func newURL(u string) *url.URL {
	res, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}

	return res
}
