package transform

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	input := map[string]interface{}{
		"type":       "json",
		"expression": "card.number",
	}

	res, err := Factory(input, nil)
	require.NoError(t, err)
	require.Equal(t, reflect.TypeOf(&Json{}), reflect.TypeOf(res))
}
