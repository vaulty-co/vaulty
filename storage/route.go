package storage

import (
	"fmt"

	"github.com/rs/xid"
	"github.com/vaulty/proxy/model"
)

func (s *redisStorage) CreateRoute(route *model.Route) error {
	route.ID = "rt" + xid.New().String()

	err := s.redisClient.Set(route.Key(), route.ID, 0).Err()
	if err != nil {
		return err
	}

	err = s.redisClient.HMSet(route.IDKey(), map[string]interface{}{
		"upstream":                 route.Upstream,
		"request_transformations":  route.RequestTransformationsJSON,
		"response_transformations": route.ResponseTransformationsJSON,
	}).Err()
	if err != nil {
		return err
	}

	listKey := fmt.Sprintf("vault:%s:routes", route.VaultID)
	err = s.redisClient.LPush(listKey, route.ID).Err()

	return err
}

func (s *redisStorage) FindRoute(vaultID string, type_ model.RouteType, method, path string) (*model.Route, error) {
	route := &model.Route{
		Type:    type_,
		Method:  method,
		Path:    path,
		VaultID: vaultID,
	}

	route.ID = s.redisClient.Get(route.Key()).Val()
	if route.ID == "" {
		return nil, ErrNoRows
	}

	return s.FindRouteByID(route.VaultID, route.ID)
}

func (s *redisStorage) FindRouteByID(vaultID, routeID string) (*model.Route, error) {
	route := &model.Route{
		ID:      routeID,
		VaultID: vaultID,
	}

	fields, err := s.redisClient.HGetAll(route.IDKey()).Result()
	if err != nil {
		return nil, err
	}

	if _, ok := fields["upstream"]; !ok {
		return nil, ErrNoRows
	}

	if len(fields) == 0 {
		return nil, ErrNoRows
	}

	route.Upstream = fields["upstream"]
	route.RequestTransformationsJSON = fields["request_transformations"]
	route.ResponseTransformationsJSON = fields["response_transformations"]

	return route, nil
}

func (s *redisStorage) ListRoutes(vaultID string) ([]*model.Route, error) {
	listKey := fmt.Sprintf("vault:%s:routes", vaultID)

	routes := []*model.Route{}

	res := s.redisClient.LRange(listKey, 0, -1)
	if res.Err() != nil {
		return nil, res.Err()
	}

	ids := res.Val()

	for _, id := range ids {
		route, err := s.FindRouteByID(vaultID, id)
		if err != nil {
			return nil, err
		}

		routes = append(routes, route)
	}

	return routes, nil
}
