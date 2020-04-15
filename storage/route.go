package storage

import (
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
		return nil, nil
	}

	fields, err := s.redisClient.HGetAll(route.IDKey()).Result()
	if err != nil {
		return nil, err
	}

	if _, ok := fields["upstream"]; !ok {
		return nil, nil
	}

	route.Upstream = fields["upstream"]
	route.RequestTransformationsJSON = fields["request_transformations"]
	route.ResponseTransformationsJSON = fields["response_transformations"]

	return route, nil
}
