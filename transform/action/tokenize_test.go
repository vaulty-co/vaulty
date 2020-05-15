package action

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/proxy/encrypt"
	"github.com/vaulty/proxy/secrets"
)

func TestTokenizeDetokenize(t *testing.T) {
	encrypter, err := encrypt.NewEncrypter("")
	require.NoError(t, err)

	secretStorage := secrets.NewEphemeralStorage(encrypter)

	plaintext := []byte("hello")

	tokenize := &Tokenize{
		secretStorage: secretStorage,
	}
	token, err := tokenize.Transform(plaintext)
	require.NoError(t, err)
	require.Contains(t, string(token), "tok")

	detokenize := &Detokenize{
		secretStorage: secretStorage,
	}
	got, err := detokenize.Transform(token)
	require.NoError(t, err)
	require.Equal(t, plaintext, got)

}
