package routing

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/encrypt"
	"github.com/vaulty/vaulty/secrets"
)

func TestLoadFromFile(t *testing.T) {
	enc, err := encrypt.NewEncrypter("")
	require.NoError(t, err)

	secretsStorage := secrets.NewEphemeralStorage(enc)

	loader := &fileLoader{
		enc:            enc,
		secretsStorage: secretsStorage,
	}

	routes, err := loader.Load("./testdata/routes.json")
	require.NoError(t, err)

	require.Equal(t, "in1", routes[0].Name)
	require.Equal(t, "in2", routes[1].Name)
	require.Equal(t, "inAll", routes[2].Name)
	require.Equal(t, "out1", routes[3].Name)
	require.Equal(t, "outAll", routes[4].Name)
}
