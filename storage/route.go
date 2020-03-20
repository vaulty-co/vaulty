package storage

import (
	"fmt"

	"github.com/vaulty/proxy/model"
)

func (s *Storage) FindRoute(vaultID, type_, method, path string) (*model.Route, error) {
	// vlt2uYBrnYkUnEF:INBOUND:POST:/records => routeID
	routeKey := fmt.Sprintf("%s:%s:%s:%s", vaultID, type_, method, path)
	routeID := s.redisClient.Get(routeKey).Val()
	if routeID == "" {
		return nil, nil
	}

	routeUpstreamKey := fmt.Sprintf("route:%s:upstream", routeID)
	routeUpstream := s.redisClient.Get(routeUpstreamKey).Val()

	return &model.Route{
		ID:       routeID,
		Upstream: routeUpstream,
	}, nil

}
