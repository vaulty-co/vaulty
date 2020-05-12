package transform

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJson(t *testing.T) {
	tr := &Json{
		Expression: "card.number",
		Action: TransformerFunc(func(body []byte) ([]byte, error) {
			return append(body, []byte(" transformed")...), nil
		}),
	}

	body := []byte(`{ "card": { "number": "4242", "cvc": "123", "exp": "10/24" } }`)
	newBody, err := tr.Transform(body)
	require.NoError(t, err)
	require.Contains(t, string(newBody), "4242 transformed")
}
