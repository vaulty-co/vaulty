package storage

import (
	"fmt"
	"log"
	"net/url"

	"github.com/vaulty/proxy/model"
)

type TestStorage struct {
}

func NewTestStorage() Storage {
	return &TestStorage{}
}

var testVaults map[string]*model.Vault = map[string]*model.Vault{}
var testRoutes map[string]*model.Route = map[string]*model.Route{}

// var testVaults map[string]*model.Vault = map[string]*model.Vault{
// 	"vlt123": &model.Vault{"vlt123", newURL("http://demo.com")},
// }

// var testRoutes map[string]*model.Route = map[string]*model.Route{
// 	"vlt123:inbound:POST:/tokenize": &model.Route{"id", newURL("http://demo.com")},
// }

// move into separate package to not add this methods and data
// to main storage package
func AddTestVault(vaultID, upstream string) {
	testVaults[vaultID] = &model.Vault{vaultID, newURL(upstream)}
}

func AddTestRoute(vaultID, type_, method, path, routeID, upstream string) {
	routeKey := fmt.Sprintf("%s:%s:%s:%s", vaultID, type_, method, path)

	testRoutes[routeKey] = &model.Route{routeID, newURL(upstream)}
}

func Reset() {
	log.Println("Reset storage")
	testRoutes = map[string]*model.Route{}
	testVaults = map[string]*model.Vault{}
}

func (s *TestStorage) FindRoute(vaultID, type_, method, path string) (*model.Route, error) {
	log.Println("Looking for route...")
	routeKey := fmt.Sprintf("%s:%s:%s:%s", vaultID, type_, method, path)

	route, ok := testRoutes[routeKey]
	if !ok {
		log.Println("Route was not found")
		return nil, nil
	}

	log.Println("Found", route)
	return route, nil
}

func (s *TestStorage) FindVault(vaultID string) (*model.Vault, error) {
	log.Print("Looking for vault...")

	vault, ok := testVaults[vaultID]
	if !ok {
		log.Println("Vault was not found")
		return nil, nil
	}

	log.Println("Found", vault)
	return vault, nil
}

func newURL(u string) *url.URL {
	res, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}

	return res
}
