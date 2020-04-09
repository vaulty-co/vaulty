package storage

import (
	"errors"

	"github.com/go-redis/redis"
	"github.com/vaulty/proxy/model"
)

var (
	ErrNoRows = errors.New("no rows found")
)

type Storage interface {
	// Vault
	CreateVault(*model.Vault) error
	FindVault(vaultID string) (*model.Vault, error)
	ListVaults() ([]*model.Vault, error)
	DeleteVault(vaultID string) error

	// Route
	CreateRoute(*model.Route) error
	FindRoute(vaultID string, type_ model.RouteType, method, path string) (*model.Route, error)
}

type redisStorage struct {
	redisClient *redis.Client
}

func NewRedisStorage(redisClient *redis.Client) Storage {
	return &redisStorage{
		redisClient: redisClient,
	}
}
