package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vaulty/proxy/model"
)

func TestWithVault(t *testing.T) {
	assert := assert.New(t)

	rs := NewRedisStorage(redisClient)
	defer redisClient.FlushAll()

	createdVault := &model.Vault{
		Upstream: "http://example.com",
	}
	err := rs.CreateVault(createdVault)
	assert.NoError(err)
	assert.NotEmpty(createdVault.ID)

	t.Run("FindVault", func(t *testing.T) {
		vault, err := rs.FindVault(createdVault.ID)
		assert.NoError(err)

		vault, err = rs.FindVault("vlt0000")

		require.Equal(t, ErrNoRows, err)
		require.Nil(t, vault)
	})

	t.Run("ListVaults", func(t *testing.T) {
		vaults, err := rs.ListVaults()
		require.NoError(t, err)
		require.Equal(t, []*model.Vault{createdVault}, vaults)
	})

	t.Run("DeleteVault", func(t *testing.T) {
		err := rs.DeleteVault(createdVault.ID)
		require.NoError(t, err)

		vault, err := rs.FindVault(createdVault.ID)
		require.Equal(t, ErrNoRows, err)
		require.Nil(t, vault)

		vaults, err := rs.ListVaults()
		require.NoError(t, err)
		require.Len(t, vaults, 0)
	})
}
