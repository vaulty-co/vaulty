package inmem

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vaulty/proxy/model"
	"github.com/vaulty/proxy/storage"
)

func TestWithVault(t *testing.T) {
	assert := assert.New(t)

	rs := NewStorage()
	defer rs.Reset()

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

		require.Equal(t, storage.ErrNoRows, err)
		require.Nil(t, vault)
	})

	t.Run("ListVaults", func(t *testing.T) {
		vaults, err := rs.ListVaults()
		require.NoError(t, err)
		require.Equal(t, []*model.Vault{createdVault}, vaults)
	})

	t.Run("DeleteVault", func(t *testing.T) {
		route := &model.Route{
			Type:     model.RouteInbound,
			Method:   http.MethodPost,
			Path:     "/tokenize",
			VaultID:  createdVault.ID,
			Upstream: "http://example.com",
		}
		err := rs.CreateRoute(route)
		require.NoError(t, err)

		err = rs.DeleteVault(createdVault.ID)
		require.NoError(t, err)

		vault, err := rs.FindVault(createdVault.ID)
		require.Equal(t, storage.ErrNoRows, err)
		require.Nil(t, vault)

		vaults, err := rs.ListVaults()
		require.NoError(t, err)
		require.Len(t, vaults, 0)

		// deletes its routes as well
		routes, err := rs.ListRoutes(createdVault.ID)
		require.NoError(t, err)
		require.Len(t, routes, 0)
	})
}
