package test_storage

import (
	"log"
	"net/url"

	"github.com/vaulty/proxy/model"
	"github.com/vaulty/proxy/storage"
)

type TestStorage struct {
}

func NewTestStorage() storage.Storage {
	return &TestStorage{}
}

// structs to store vaults and routes
type vaultRouteKey struct {
	vaultID, routeID string
}

var testVaults map[string]*model.Vault = map[string]*model.Vault{}
var testRoutes map[string]*model.Route = map[string]*model.Route{}
var testVaultRoutesIDs map[vaultRouteKey]*model.Route = map[vaultRouteKey]*model.Route{}

func Reset() {
	testRoutes = map[string]*model.Route{}
	testVaultRoutesIDs = map[vaultRouteKey]*model.Route{}
	testVaults = map[string]*model.Vault{}
}

func newURL(u string) *url.URL {
	res, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}

	return res
}
