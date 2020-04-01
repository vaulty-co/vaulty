package storage

import (
	"github.com/vaulty/proxy/model"
)

type Storage interface {
	FindRoute(vaultID, type_, method, path string) (*model.Route, error)
	FindVault(vaultID string) (*model.Vault, error)
}
