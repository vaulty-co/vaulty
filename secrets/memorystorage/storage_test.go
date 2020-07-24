package memorystorage

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/config"
	"github.com/vaulty/vaulty/encryption/noneenc"
	"github.com/vaulty/vaulty/secrets"
)

func TestFactory(t *testing.T) {
	enc, err := noneenc.Factory(&config.Config{})
	require.NoError(t, err)

	store, err := Factory(&secrets.Config{
		Encrypter: enc,
	})
	require.NoError(t, err)
	require.Implements(t, (*secrets.Storage)(nil), store)
}

func TestMemoryStorage(t *testing.T) {
	encrypter := noneenc.New()

	storage := New(&Params{
		Encrypter: encrypter,
	})

	value := []byte("1")
	err := storage.Set("one", value)
	require.NoError(t, err)

	got, err := storage.Get("one")
	require.NoError(t, err)
	require.Equal(t, value, got)

}
