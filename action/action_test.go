package action

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	input := map[string]interface{}{
		"type": "tokenize",
	}

	res, err := Factory(input, &Options{})
	require.NoError(t, err)
	require.Equal(t, reflect.TypeOf(&Tokenize{}), reflect.TypeOf(res))
}
