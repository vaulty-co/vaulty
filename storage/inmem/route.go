package inmem

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/rs/xid"
	"github.com/vaulty/proxy/model"
	"github.com/vaulty/proxy/storage"
)

func (s *inmemStorage) CreateRoute(route *model.Route) error {
	route.ID = "rt" + xid.New().String()

	fmt.Println("(create) route key", route.Key())
	s.routes[route.Key()] = route
	s.vaultRoutesIDs[vaultRouteKey{vaultID: route.VaultID, routeID: route.ID}] = route

	return nil
}

func (s *inmemStorage) FindRoute(vaultID string, type_ model.RouteType, req *http.Request) (*model.Route, error) {
	var target string

	if type_ == model.RouteInbound {
		target = req.URL.Path
	} else {
		matchingURL := &url.URL{}
		matchingURL.Host = req.URL.Host
		matchingURL.Scheme = req.URL.Scheme
		matchingURL.Path = req.URL.Path
		target = matchingURL.String()
	}

	routeKey := fmt.Sprintf("%s:%s:%s:%s", vaultID, type_, req.Method, target)

	fmt.Println("(find) route key", routeKey)

	route, ok := s.routes[routeKey]
	if !ok {
		return nil, storage.ErrNoRows
	}

	return route, nil
}

func (s *inmemStorage) FindRouteByID(vaultID, routeID string) (*model.Route, error) {
	route, ok := s.vaultRoutesIDs[vaultRouteKey{vaultID: vaultID, routeID: routeID}]
	if !ok {
		// route was not found
		return nil, storage.ErrNoRows
	}

	return route, nil
}

func (s *inmemStorage) ListRoutes(vaultID string) ([]*model.Route, error) {
	routes := []*model.Route{}

	for _, r := range s.vaultRoutesIDs {
		if r.VaultID != vaultID {
			continue
		}

		tmp := *r
		route := &tmp

		routes = append(routes, route)
	}

	return routes, nil
}

func (s *inmemStorage) DeleteRoute(vaultID, routeID string) error {
	route, err := s.FindRouteByID(vaultID, routeID)
	if err != nil {
		return err
	}

	delete(s.routes, route.Key())
	delete(s.vaultRoutesIDs, vaultRouteKey{vaultID: route.VaultID, routeID: route.ID})
	return nil
}

func (s *inmemStorage) DeleteRoutes(vaultID string) error {
	for key, route := range s.vaultRoutesIDs {
		if route.VaultID != vaultID {
			continue
		}

		delete(s.routes, route.Key())
		delete(s.vaultRoutesIDs, key)

	}

	return nil
}
