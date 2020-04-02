package storage

import (
	"github.com/vaulty/proxy/model"
)

type Storage interface {
	CreateRoute(*model.Route) error
	FindRoute(vaultID string, type_ model.RouteType, method, path string) (*model.Route, error)
	FindVault(vaultID string) (*model.Vault, error)
}
