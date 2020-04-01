package test_storage

import (
	"fmt"
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

var testVaults map[string]*model.Vault = map[string]*model.Vault{}
var testRoutes map[string]*model.Route = map[string]*model.Route{}

func AddTestVault(vaultID, upstream string) {
	testVaults[vaultID] = &model.Vault{vaultID, newURL(upstream)}
}

func AddTestRoute(vaultID, type_, method, path, routeID, upstream string) {
	routeKey := fmt.Sprintf("%s:%s:%s:%s", vaultID, type_, method, path)

	testRoutes[routeKey] = &model.Route{routeID, newURL(upstream)}
}

func Reset() {
	testRoutes = map[string]*model.Route{}
	testVaults = map[string]*model.Vault{}
}

func (s *TestStorage) FindRoute(vaultID, type_, method, path string) (*model.Route, error) {
	routeKey := fmt.Sprintf("%s:%s:%s:%s", vaultID, type_, method, path)

	route, ok := testRoutes[routeKey]
	if !ok {
		// route was not found
		return nil, nil
	}

	return route, nil
}

func (s *TestStorage) FindVault(vaultID string) (*model.Vault, error) {
	vault, ok := testVaults[vaultID]
	if !ok {
		// vault was not found
		return nil, nil
	}

	return vault, nil
}

func newURL(u string) *url.URL {
	res, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}

	return res
}
