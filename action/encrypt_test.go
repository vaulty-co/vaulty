package action

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/encryption/noneenc"
)

func TestEncryptDecrypt(t *testing.T) {
	encrypter := noneenc.New()

	plaintext := []byte("hello")

	encrypt := &Encrypt{
		enc: encrypter,
	}
	ciphertext, err := encrypt.Transform(plaintext)
	require.NoError(t, err)
	require.NotEqual(t, ciphertext, plaintext)

	decrypt := &Decrypt{
		enc: encrypter,
	}
	decrypted, err := decrypt.Transform(ciphertext)
	require.NoError(t, err)
	require.Equal(t, plaintext, decrypted)
}
