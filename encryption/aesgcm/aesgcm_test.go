package aesgcm

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/config"
	"github.com/vaulty/vaulty/encryption"
)

func TestFactory(t *testing.T) {
	_, err := Factory(&config.Config{
		Encryption: &config.Encryption{
			Key: "123",
		},
	})
	require.EqualError(t, err, "invalid key length: 3. Should be 32 bytes")

	enc, err := Factory(&config.Config{
		Encryption: &config.Encryption{
			Key: "776f726420746f206120736563726574",
		},
	})
	require.NoError(t, err)
	require.Implements(t, (*encryption.Encrypter)(nil), enc)
}

func TestAesGCM(t *testing.T) {
	enc, err := NewEncrypter([]byte("776f726420746f206120736563726574"))
	require.NoError(t, err)

	ciphertext, err := enc.Encrypt([]byte("hello"))
	require.NoError(t, err)
	require.NotEqual(t, []byte("hello"), ciphertext)

	plaintext, err := enc.Decrypt(ciphertext)
	require.NoError(t, err)
	require.Equal(t, []byte("hello"), plaintext)
}
