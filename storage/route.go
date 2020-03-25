package storage

import (
	"fmt"
	"net/url"

	"github.com/vaulty/proxy/model"
)

func (s *Storage) FindRoute(vaultID, type_, method, path string) (*model.Route, error) {
	// vlt2uYBrnYkUnEF:INBOUND:POST:/records => routeID
	routeKey := fmt.Sprintf("%s:%s:%s:%s", vaultID, type_, method, path)
	routeID := s.redisClient.Get(routeKey).Val()
	if routeID == "" {
		return nil, nil
	}

	upstreamKey := fmt.Sprintf("route:%s:upstream", routeID)
	upstream := s.redisClient.Get(upstreamKey).Val()
	upstreamURL, err := url.Parse(upstream)
	if err != nil {
		return nil, err
	}

	return &model.Route{
		ID:          routeID,
		UpstreamURL: upstreamURL,
	}, nil

}
