package storage

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vaulty/proxy/model"
)

func TestFindRoute(t *testing.T) {
	assert := assert.New(t)

	rs := NewRedisStorage(redisClient)

	err := rs.CreateRoute(&model.Route{
		ID:       "rt1",
		Type:     model.RouteInbound,
		Method:   http.MethodPost,
		Path:     "/tokenize",
		VaultID:  "vlt1",
		Upstream: "http://example.com",
	})
	assert.NoError(err)

	route, err := rs.FindRoute("vlt1", model.RouteInbound, http.MethodPost, "/tokenize")

	assert.NoError(err)
	assert.Equal(route.ID, "rt1")

	route, err = rs.FindRoute("vlt1", model.RouteInbound, http.MethodPost, "/nothing")
	assert.NoError(err)
	assert.Nil(route)
}

func TestFindVault(t *testing.T) {
	assert := assert.New(t)

	rs := NewRedisStorage(redisClient)

	err := rs.CreateVault(&model.Vault{
		ID:       "vlt1",
		Upstream: "http://example.com",
	})
	assert.NoError(err)

	vault, err := rs.FindVault("vlt1")
	assert.NoError(err)
	assert.Equal(vault.ID, "vlt1")

	vault, err = rs.FindVault("vlt0")
	assert.NoError(err)
	assert.Nil(vault)
}
