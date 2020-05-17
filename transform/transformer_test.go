package transform

import (
	"encoding/json"
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

// helper method that is used by other transformations tests
func transformerFromJSON(rawJson []byte, action Transformer) (Transformer, error) {
	var input map[string]interface{}

	err := json.Unmarshal(rawJson, &input)
	if err != nil {
		return nil, err
	}

	transformation, err := Factory(input, action)
	if err != nil {
		return nil, err
	}

	return transformation, nil
}
