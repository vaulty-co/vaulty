package action

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/encryption/noneenc"
	"github.com/vaulty/vaulty/secrets/memorystorage"
	"crypto/sha256"
	"encoding/hex"
)

func TestTokenizeAndHashWithFormat(t *testing.T) {
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
