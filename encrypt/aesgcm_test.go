package encrypt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAesGCM(t *testing.T) {
	enc, err := NewAesGcm("776f726420746f206120736563726574")
	require.NoError(t, err)

	ciphertext, err := enc.Encrypt([]byte("hello"))
	require.NoError(t, err)
	require.NotEqual(t, []byte("hello"), ciphertext)

	plaintext, err := enc.Decrypt(ciphertext)
	require.NoError(t, err)
	require.Equal(t, []byte("hello"), plaintext)
}
