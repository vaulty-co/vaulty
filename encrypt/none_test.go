package encrypt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNone(t *testing.T) {
	enc := &None{}

	ciphertext, err := enc.Encrypt([]byte("hello"))
	require.NoError(t, err)
	require.NotEqual(t, []byte("hello"), ciphertext)

	fmt.Println(string(ciphertext))

	plaintext, err := enc.Decrypt(ciphertext)
	require.NoError(t, err)
	require.Equal(t, []byte("hello"), plaintext)
}
