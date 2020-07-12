package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Cleanup(func() {
		os.Unsetenv("VAULTY_ENCRYPTION_KEY")
		os.Unsetenv("VAULTY_ENCRYPTION_TYPE")
	})

	os.Setenv("VAULTY_ENCRYPTION_KEY", "1234567890")
	os.Setenv("VAULTY_ENCRYPTION_TYPE", "aesgcm")

	cnf := &Config{}

	err := cnf.FromEnvironment()
	require.NoError(t, err)

	require.Equal(t, "aesgcm", cnf.Encryption.Type)
	require.Equal(t, "1234567890", cnf.Encryption.Key)
}
