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

	t.Run("Test transformation validation", func(t *testing.T) {
		tr := &Json{
			Expression: "users.#.email",
		}

		require.NoError(t, tr.Validate())

		tr = &Json{
			Expression: "one.#.users.#.email",
		}

		require.EqualError(t, tr.Validate(), "Nested arrays are not supported, but used in the expression: one.#.users.#.email")
	})

	t.Run("Test transformation of array", func(t *testing.T) {
		tr := &Json{
			Expression: "users.#.email",
			Action: TransformerFunc(func(body []byte) ([]byte, error) {
				return append(body, '+', body[0]), nil
			}),
		}

		body := []byte(`
{
	"users":[
		{
			"id":1,
			"email":"a@mail.com"
		},
		{
			"id":2,
			"email":"b@mail.com"
		},
		{
			"id":3,
			"email":"c@mail.com"
		}
	]

`)

		want := []byte(`
{
	"users":[
		{
			"id":1,
			"email":"a@mail.com+a"
		},
		{
			"id":2,
			"email":"b@mail.com+b"
		},
		{
			"id":3,
			"email":"c@mail.com+c"
		}
	]

`)
		newBody, err := tr.Transform(body)
		require.NoError(t, err)
		require.Equal(t, string(want), string(newBody))
	})
}
