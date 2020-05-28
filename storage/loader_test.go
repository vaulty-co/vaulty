package storage_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/storage"
	"github.com/vaulty/vaulty/storage/inmem"
	"github.com/vaulty/vaulty/transform/action"
)

func TestLoadFromFile(t *testing.T) {
	st := inmem.NewStorage()

	err := storage.LoadFromFile("./test-fixture/routes.json", &storage.LoaderOptions{
		Storage:       st,
		ActionOptions: &action.Options{},
	})
	require.NoError(t, err)

	vaults, err := st.ListVaults()
	require.NoError(t, err)
	require.Len(t, vaults, 1)

	routes, err := st.ListRoutes(vaults[0].ID)
	require.NoError(t, err)
	require.Len(t, routes, 2)
}
