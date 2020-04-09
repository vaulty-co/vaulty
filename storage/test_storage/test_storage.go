package test_storage

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/rs/xid"
	"github.com/vaulty/proxy/model"
	"github.com/vaulty/proxy/storage"
)

type TestStorage struct {
}

func NewTestStorage() storage.Storage {
	return &TestStorage{}
}

var testVaults map[string]*model.Vault = map[string]*model.Vault{}
var testRoutes map[string]*model.Route = map[string]*model.Route{}

func Reset() {
	testRoutes = map[string]*model.Route{}
	testVaults = map[string]*model.Vault{}
}

func (s *TestStorage) CreateRoute(route *model.Route) error {
	testRoutes[route.Key()] = route

	return nil
}

func (s *TestStorage) CreateVault(vault *model.Vault) error {
	vault.ID = "vlt" + xid.New().String()

	testVaults[vault.ID] = vault

	return nil
}

func (s *TestStorage) ListVaults() ([]*model.Vault, error) {
	vaults := []*model.Vault{}

	for _, v := range testVaults {
		vault := &model.Vault{}
		vault.ID = v.ID
		vault.Upstream = v.Upstream

		vaults = append(vaults, vault)
	}

	return vaults, nil
}

func (s *TestStorage) FindRoute(vaultID string, type_ model.RouteType, method, path string) (*model.Route, error) {
	routeKey := fmt.Sprintf("%s:%s:%s:%s", vaultID, type_, method, path)

	route, ok := testRoutes[routeKey]
	if !ok {
		// route was not found
		return nil, nil
	}

	return route, nil
}

func (s *TestStorage) FindVault(vaultID string) (*model.Vault, error) {
	if vaultID == "vltError" {
		return nil, errors.New("Test error")
	}

	vault, ok := testVaults[vaultID]
	if !ok {
		// vault was not found
		return nil, storage.ErrNoRows
	}

	return vault, nil
}

func (s *TestStorage) DeleteVault(vaultID string) error {
	delete(testVaults, vaultID)

	return nil
}

func newURL(u string) *url.URL {
	res, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}

	return res
}
