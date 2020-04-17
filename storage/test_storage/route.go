package test_storage

import (
	"fmt"

	"github.com/rs/xid"
	"github.com/vaulty/proxy/model"
	"github.com/vaulty/proxy/storage"
)

func (s *TestStorage) CreateRoute(route *model.Route) error {
	route.ID = "rt" + xid.New().String()

	testRoutes[route.Key()] = route
	testVaultRoutesIDs[vaultRouteKey{vaultID: route.VaultID, routeID: route.ID}] = route

	return nil
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

func (s *TestStorage) FindRouteByID(vaultID, routeID string) (*model.Route, error) {
	route, ok := testVaultRoutesIDs[vaultRouteKey{vaultID: vaultID, routeID: routeID}]
	if !ok {
		// route was not found
		return nil, storage.ErrNoRows
	}

	return route, nil
}

func (s *TestStorage) ListRoutes(vaultID string) ([]*model.Route, error) {
	routes := []*model.Route{}

	for _, r := range testVaultRoutesIDs {
		if r.VaultID != vaultID {
			continue
		}

		tmp := *r
		route := &tmp

		routes = append(routes, route)
	}

	return routes, nil
}

func (s *TestStorage) DeleteRoute(vaultID, routeID string) error {
	route, err := s.FindRouteByID(vaultID, routeID)
	if err != nil {
		return err
	}

	delete(testRoutes, route.Key())
	delete(testVaultRoutesIDs, vaultRouteKey{vaultID: route.VaultID, routeID: route.ID})
	return nil
}

func (s *TestStorage) DeleteRoutes(vaultID string) error {
	for key, route := range testVaultRoutesIDs {
		if route.VaultID != vaultID {
			continue
		}

		delete(testRoutes, route.Key())
		delete(testVaultRoutesIDs, key)

	}

	return nil
}
