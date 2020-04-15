package storage

import (
	"net/http"
	"os"
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vaulty/proxy/core"
	"github.com/vaulty/proxy/model"
)

var redisClient *redis.Client

func TestMain(m *testing.M) {
	core.LoadConfig("../config/test.yml")

	redisOptions, err := redis.ParseURL(core.Config.Redis.URL)
	if err != nil {
		panic(err)
	}

	redisClient = redis.NewClient(redisOptions)
	redisClient.FlushAll()
	exitCode := m.Run()
	redisClient.FlushAll()
	os.Exit(exitCode)
}

func TestFindRoute(t *testing.T) {
	assert := require.New(t)

	rs := NewRedisStorage(redisClient)
	defer redisClient.FlushAll()

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
	defer redisClient.FlushAll()

	vault := &model.Vault{
		Upstream: "http://example.com",
	}
	err := rs.CreateVault(vault)
	assert.NoError(err)
	assert.NotEmpty(vault.ID)

	vault, err = rs.FindVault(vault.ID)
	assert.NoError(err)

	vault, err = rs.FindVault("vlt0000")

	require.Equal(t, ErrNoRows, err)
	require.Nil(t, vault)
}

func TestListVaults(t *testing.T) {
	rs := NewRedisStorage(redisClient)
	defer redisClient.FlushAll()

	vault := &model.Vault{
		Upstream: "http://example.com",
	}
	err := rs.CreateVault(vault)
	require.NoError(t, err)

	vaults, err := rs.ListVaults()
	require.NoError(t, err)
	require.Equal(t, []*model.Vault{vault}, vaults)
}

func TestDeleteVault(t *testing.T) {
	rs := NewRedisStorage(redisClient)
	defer redisClient.FlushAll()

	vault := &model.Vault{
		Upstream: "http://example.com",
	}
	err := rs.CreateVault(vault)
	require.NoError(t, err)

	err = rs.DeleteVault(vault.ID)
	require.NoError(t, err)

	vault, err = rs.FindVault(vault.ID)
	require.Equal(t, ErrNoRows, err)
	require.Nil(t, vault)

	vaults, err := rs.ListVaults()
	require.NoError(t, err)
	require.Len(t, vaults, 0)

}
