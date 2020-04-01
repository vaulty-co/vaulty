package storage

import (
	"fmt"
	"net/url"

	"github.com/go-redis/redis"
	"github.com/vaulty/proxy/model"
)

type RedisStorage struct {
	redisClient *redis.Client
}

func NewRedisStorage(redisClient *redis.Client) Storage {
	return &RedisStorage{
		redisClient: redisClient,
	}
}

func (s *RedisStorage) FindRoute(vaultID, type_, method, path string) (*model.Route, error) {
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

func (s *RedisStorage) FindVault(vaultID string) (*model.Vault, error) {
	upstreamKey := fmt.Sprintf("vault:%s:upstream", vaultID)
	upstream := s.redisClient.Get(upstreamKey).Val()
	if upstream == "" {
		return nil, nil
	}

	upstreamURL, err := url.Parse(upstream)
	if err != nil {
		return nil, err
	}

	return &model.Vault{
		ID:          vaultID,
		UpstreamURL: upstreamURL,
	}, nil
}
