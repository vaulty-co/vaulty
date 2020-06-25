package action

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	action := &Hash{}
	result, err := action.Transform([]byte("hello"))
	require.NoError(t, err)
	require.Equal(t, "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824", string(result))

	actionWithSalt := &Hash{
		salt: "salt",
	}
	resultWithSalt, err := actionWithSalt.Transform([]byte("hello"))
	require.NoError(t, err)
	require.NotEqual(t, result, resultWithSalt)
	require.Equal(t, "87daba3fe263b34c335a0ee3b28ffec4d159aad6542502eaf551dc7b9128c267", string(resultWithSalt))
}
