package secrets

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/proxy/encrypt"
)

func TestEphemeral(t *testing.T) {
	encrypter, err := encrypt.NewEncrypter("")
	require.NoError(t, err)

	storage := NewEphemeralStorage(encrypter)

	value := []byte("1")
	err = storage.Set("one", value)
	require.NoError(t, err)

	got, err := storage.Get("one")
	require.NoError(t, err)
	require.Equal(t, value, got)

}
