package storage

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vaulty/proxy/model"
)

func TestFindRoute(t *testing.T) {
	assert := require.New(t)

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
	assert.NotNil(route)

	assert.Equal(route.ID, "rt1")
	assert.Equal("http://example.com", route.Upstream)

	route, err = rs.FindRoute("vlt1", model.RouteInbound, http.MethodPost, "/nothing")
	assert.NoError(err)
	assert.Nil(route)
}

func TestFindVault(t *testing.T) {
	assert := assert.New(t)

	rs := NewRedisStorage(redisClient)

	vault := &model.Vault{
		Upstream: "http://example.com",
	}
	err := rs.CreateVault(vault)
	assert.NoError(err)
	assert.NotEmpty(vault.ID)

	vault, err = rs.FindVault(vault.ID)
	assert.NoError(err)

	vault, err = rs.FindVault("vlt0000")
	assert.NoError(err)
	assert.Nil(vault)
}
