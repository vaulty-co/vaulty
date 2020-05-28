package storage

import (
	"errors"
	"net/http"

	"github.com/vaulty/vaulty/model"
)

var (
	ErrNoRows = errors.New("no rows found")
)

type Storage interface {
	Reset()

	// Vault
	CreateVault(*model.Vault) error
	FindVault(vaultID string) (*model.Vault, error)
	ListVaults() ([]*model.Vault, error)
	DeleteVault(vaultID string) error
	UpdateVault(*model.Vault) error

	// Route
	CreateRoute(*model.Route) error
	FindRoute(vaultID string, type_ model.RouteType, req *http.Request) (*model.Route, error)
	FindRouteByID(vaultID, routeID string) (*model.Route, error)
	ListRoutes(vaultID string) ([]*model.Route, error)
	DeleteRoute(vaultID, routeID string) error
	DeleteRoutes(vaultID string) error
}
