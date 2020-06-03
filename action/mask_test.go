package action

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMask(t *testing.T) {
	action := &Mask{}
	result, err := action.Transform([]byte("hello"))
	require.NoError(t, err)
	require.Equal(t, "*****", string(result))

	action2 := &Mask{
		Symbol: []byte("x"),
	}
	result, err = action2.Transform([]byte("hello"))
	require.NoError(t, err)
	require.Equal(t, "xxxxx", string(result))
}
