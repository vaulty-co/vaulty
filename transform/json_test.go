package transform

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJson(t *testing.T) {
	t.Run("Test transformation", func(t *testing.T) {
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
	})

	t.Run("Test transformation with invalid json", func(t *testing.T) {
		tr := &Json{
			Expression: "card.number",
			Action: TransformerFunc(func(body []byte) ([]byte, error) {
				return append(body, []byte(" transformed")...), nil
			}),
		}

		body := []byte(`not valid json`)
		newBody, err := tr.Transform(body)
		require.NoError(t, err)
		require.Equal(t, newBody, body)
	})
}
