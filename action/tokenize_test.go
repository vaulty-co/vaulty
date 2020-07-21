package action

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/encryption/noneenc"
	"github.com/vaulty/vaulty/secrets/memorystorage"
)

func TestTokenizeDetokenize(t *testing.T) {
	encrypter := noneenc.New()

	secretsStorage := memorystorage.New(&memorystorage.Params{
		Encrypter: encrypter,
	})

	plaintext := []byte("hello")

	tokenize := &Tokenize{
		secretsStorage: secretsStorage,
	}
	token, err := tokenize.Transform(plaintext)
	require.NoError(t, err)
	require.Contains(t, string(token), "tok")

	detokenize := &Detokenize{
		secretsStorage: secretsStorage,
	}
	got, err := detokenize.Transform(token)
	require.NoError(t, err)
	require.Equal(t, plaintext, got)
}

func TestTokenizeWithFormat(t *testing.T) {
	encrypter := noneenc.New()

	secretsStorage := memorystorage.New(&memorystorage.Params{
		Encrypter: encrypter,
	})

	plaintext := []byte("hello")

	tokenize := &Tokenize{
		secretsStorage: secretsStorage,
		Format:         "email",
	}
	token, err := tokenize.Transform(plaintext)
	require.NoError(t, err)
	require.Contains(t, string(token), "tok")
}
