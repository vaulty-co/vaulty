package json

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/action"
)

func TestJson(t *testing.T) {
	t.Run("Test building transformation from JSON", func(t *testing.T) {
		rawJson := []byte(`
		{
			"type":"json",
			"expression":"user.email"
		}
		`)

		var input map[string]interface{}
		err := json.Unmarshal(rawJson, &input)

		fakeAction := action.ActionFunc(func(body []byte) ([]byte, error) {
			return nil, nil
		})
		transformation, err := Factory(input, fakeAction)
		require.NoError(t, err)
		require.NotNil(t, transformation)
	})

	t.Run("Test transformation validation", func(t *testing.T) {
		_, err := NewTransformation(&Params{Expression: "users.#.email"})
		require.NoError(t, err)

		_, err = NewTransformation(&Params{Expression: "one.#.users.#.email"})
		require.EqualError(t, err, "Nested arrays are not supported, but used in the expression: one.#.users.#.email")
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

	t.Run("Test transformation with multiple paths", func(t *testing.T) {
		tr, err := NewTransformation(&Params{
			Expression: "card.number, card.cvc",
			Action: action.ActionFunc(func(body []byte) ([]byte, error) {
				return append(body, '+', body[0]), nil
			}),
		})
		require.NoError(t, err)

		body := []byte(`{"card": {"number":"1234", "cvc":"5678"}}`)

		want := []byte(`{"card": {"number":"1234+1", "cvc":"5678+5"}}`)
		newBody, err := tr.Transform(body)
		require.NoError(t, err)
		require.Equal(t, string(want), string(newBody))
	})

	t.Run("Test request transformation", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/request", strings.NewReader(`{"name": "John"}`))
		req.Header.Set("Content-Type", "application/json")

		tr, err := NewTransformation(&Params{
			Expression: "name",
			Action: action.ActionFunc(func(body []byte) ([]byte, error) {
				return []byte("transformed"), nil
			}),
		})
		require.NoError(t, err)

		req, err = tr.TransformRequest(req)
		require.NoError(t, err)
		newBody, err := ioutil.ReadAll(req.Body)
		require.Equal(t, `{"name": "transformed"}`, string(newBody))
	})

	t.Run("Test request transformation with unsupported content type", func(t *testing.T) {
		body := ioutil.NopCloser(strings.NewReader(`{"name": "John"}`))
		req, _ := http.NewRequest("POST", "/request", body)
		req.Header.Set("Content-Type", "text/plain")

		tr, err := NewTransformation(&Params{
			Expression: "name",
		})
		require.NoError(t, err)

		newReq, err := tr.TransformRequest(req)
		require.NoError(t, err)
		require.Equal(t, req, newReq)
		require.Equal(t, body, newReq.Body)
	})

	t.Run("Test response transformation", func(t *testing.T) {
		res := &http.Response{
			Body:   ioutil.NopCloser(strings.NewReader(`{"name": "John"}`)),
			Header: http.Header{"Content-Type": {"application/json"}},
		}

		tr, err := NewTransformation(&Params{
			Expression: "name",
			Action: action.ActionFunc(func(body []byte) ([]byte, error) {
				return []byte("transformed"), nil
			}),
		})
		require.NoError(t, err)

		res, err = tr.TransformResponse(res)
		require.NoError(t, err)
		newBody, err := ioutil.ReadAll(res.Body)
		require.Equal(t, `{"name": "transformed"}`, string(newBody))
	})

	t.Run("Test response transformation with unsupported content type", func(t *testing.T) {
		body := ioutil.NopCloser(strings.NewReader(`{"name": "John"}`))
		res := &http.Response{
			Body:   body,
			Header: http.Header{"Content-Type": {"text/plain"}},
		}

		tr, err := NewTransformation(&Params{
			Expression: "name",
		})
		require.NoError(t, err)

		newRes, err := tr.TransformResponse(res)
		require.NoError(t, err)
		require.Equal(t, res, newRes)
		require.Equal(t, body, newRes.Body)
	})
}
