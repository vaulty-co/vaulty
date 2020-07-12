package noneenc

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/config"
	"github.com/vaulty/vaulty/encryption"
)

func TestFactory(t *testing.T) {
	enc, err := Factory(&config.Config{})
	require.NoError(t, err)
	require.Implements(t, (*encryption.Encrypter)(nil), enc)
}

func TestNone(t *testing.T) {
	enc := &None{}

	ciphertext, err := enc.Encrypt([]byte("hello"))
	require.NoError(t, err)
	require.NotEqual(t, []byte("hello"), ciphertext)

	plaintext, err := enc.Decrypt(ciphertext)
	require.NoError(t, err)
	require.Equal(t, []byte("hello"), plaintext)
}
