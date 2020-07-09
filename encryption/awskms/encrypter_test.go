package awskms

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncrypter(t *testing.T) {
	if os.Getenv("KMS_KEY_ID") == "" {
		t.Skip("Skipping AWS KMS test as KMS_KEY_ID is not set in env")
	}

	enc, err := NewEncrypter(&Params{
		KeyID:  os.Getenv("KMS_KEY_ID"),
		Region: "eu-central-1",
	})

	message, err := enc.Encrypt([]byte("hello world!"))
	require.NoError(t, err)

	plaintext, err := enc.Decrypt(message)
	require.NoError(t, err)

	require.Equal(t, "hello world!", string(plaintext))
}
