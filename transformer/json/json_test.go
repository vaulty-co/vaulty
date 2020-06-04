package json

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/action"
)

func TestJson(t *testing.T) {
	t.Run("Test building transformer from JSON", func(t *testing.T) {
		rawJson := []byte(`
		{
			"type":"json",
			"expression":"user.email"
		}
		`)

		var input map[string]interface{}
		err := json.Unmarshal(rawJson, &input)

		fakeAction := action.ActionFunc(func(body []byte) ([]byte, error) {
			require.Equal(t, []byte("john@example.com"), body)
			return body, nil
		})

		transformation, err := Factory(input, fakeAction)
		require.NoError(t, err)

		body := bytes.NewBufferString(`{ "user": { "name": "John", "email": "john@example.com" } }`)
		req := httptest.NewRequest("POST", "/url", body)
		_, err = transformation.TransformRequest(req)
		require.NoError(t, err)

	})

	t.Run("Test transformation", func(t *testing.T) {
		tr, err := NewTransformation(&Params{
			Expression: "card.number",
			Action: action.ActionFunc(func(body []byte) ([]byte, error) {
				require.Equal(t, []byte("4242"), body)
				return body, nil
			}),
		})
		require.NoError(t, err)

		body := []byte(`{ "card": { "number": "4242", "cvc": "123", "exp": "10/24" } }`)
		_, err = tr.Transform(body)
		require.NoError(t, err)
	})

	t.Run("Test transformation with invalid json", func(t *testing.T) {
		tr, err := NewTransformation(&Params{
			Expression: "card.number",
			Action: action.ActionFunc(func(body []byte) ([]byte, error) {
				require.Fail(t, "Should not be called")
				return nil, nil
			}),
		})
		require.NoError(t, err)

		body := []byte(`not valid json`)
		newBody, err := tr.Transform(body)
		require.NoError(t, err)
		require.Equal(t, newBody, body)
	})

	t.Run("Test factory validation", func(t *testing.T) {
		_, err := NewTransformation(&Params{Expression: "users.#.email"})
		require.NoError(t, err)

		_, err = NewTransformation(&Params{Expression: "one.#.users.#.email"})
		require.EqualError(t, err, "Nested arrays are not supported, but used in the expression: one.#.users.#.email")
	})

	t.Run("Test transformation of array", func(t *testing.T) {
		tr, err := NewTransformation(&Params{
			Expression: "users.#.email",
			Action: action.ActionFunc(func(body []byte) ([]byte, error) {
				return append(body, '+', body[0]), nil
			}),
		})
		require.NoError(t, err)

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
		},
		{
			"id":4,
			"email": true
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
		},
		{
			"id":4,
			"email": true
		}
	]

	`)
		newBody, err := tr.Transform(body)
		require.NoError(t, err)
		require.Equal(t, string(want), string(newBody))
	})
}
