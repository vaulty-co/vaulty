package storage

import (
	"github.com/go-redis/redis"
	"github.com/vaulty/proxy/model"
)

type Storage interface {
	CreateRoute(*model.Route) error
	CreateVault(*model.Vault) error
	ListVaults() ([]*model.Vault, error)
	FindRoute(vaultID string, type_ model.RouteType, method, path string) (*model.Route, error)
	FindVault(vaultID string) (*model.Vault, error)
}

type redisStorage struct {
	redisClient *redis.Client
}

func NewRedisStorage(redisClient *redis.Client) Storage {
	return &redisStorage{
		redisClient: redisClient,
	}
}
