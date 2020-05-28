package request

import (
	"context"

	"github.com/vaulty/vaulty/model"
)

type key int

const (
	vaultKey key = iota
	routeKey
)

func WithVault(parent context.Context, vault *model.Vault) context.Context {
	return context.WithValue(parent, vaultKey, vault)
}

func WithRoute(parent context.Context, route *model.Route) context.Context {
	return context.WithValue(parent, routeKey, route)
}

func VaultFrom(ctx context.Context) *model.Vault {
	vault, ok := ctx.Value(vaultKey).(*model.Vault)
	if !ok {
		return nil
	}

	return vault
}

func RouteFrom(ctx context.Context) *model.Route {
	route, ok := ctx.Value(routeKey).(*model.Route)
	if !ok {
		return nil
	}

	return route
}
