package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vaulty/proxy/model"
)

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
