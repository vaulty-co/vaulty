package request

import (
	"context"

	"github.com/vaulty/proxy/model"
)

type key int

const (
	vaultKey key = iota
	routeKey
)

func WithVault(parent context.Context, vault *model.Vault) context.Context {
	return context.WithValue(parent, vaultKey, vault)
}

func VaultFrom(ctx context.Context) (*model.Vault, bool) {
	vault, ok := ctx.Value(vaultKey).(*model.Vault)
	return vault, ok
}
