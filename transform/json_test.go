package transform

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJson(t *testing.T) {
	t.Run("Test building transformer from JSON", func(t *testing.T) {
		rawJson := []byte(`
{
	"type":"json",
	"expression":"user.email"
}
`)

		fakeAction := TransformerFunc(func(body []byte) ([]byte, error) {
			require.Equal(t, []byte("john@example.com"), body)
			return body, nil
		})
		transformation, err := transformerFromJSON(rawJson, fakeAction)
		require.NoError(t, err)

		body := []byte(`{ "user": { "name": "John", "email": "john@example.com" } }`)

		body, err = transformation.Transform(body)
	})

	t.Run("Test transformation", func(t *testing.T) {
		tr := &Json{
			Expression: "card.number",
			Action: TransformerFunc(func(body []byte) ([]byte, error) {
				require.Equal(t, []byte("4242"), body)
				return body, nil
			}),
		}

		body := []byte(`{ "card": { "number": "4242", "cvc": "123", "exp": "10/24" } }`)
		_, err := tr.Transform(body)
		require.NoError(t, err)
	})

	t.Run("Test transformation with invalid json", func(t *testing.T) {
		tr := &Json{
			Expression: "card.number",
			Action: TransformerFunc(func(body []byte) ([]byte, error) {
				require.Fail(t, "Should not be called")
				return nil, nil
			}),
		}

		body := []byte(`not valid json`)
		newBody, err := tr.Transform(body)
		require.NoError(t, err)
		require.Equal(t, newBody, body)
	})
}
