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

func (s *RedisStorage) CreateRoute(route *model.Route) error {
	err := s.redisClient.Set(route.Key(), route.ID, 0).Err()
	return err
}

func (s *RedisStorage) FindRoute(vaultID string, type_ model.RouteType, method, path string) (*model.Route, error) {
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

	route.Upstream = s.redisClient.Get(route.UpstreamKey()).Val()

	return route, nil
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
