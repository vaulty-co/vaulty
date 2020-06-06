package regexp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/action"
)

func TestRegexp(t *testing.T) {
	t.Run("Test building transformer from JSON", func(t *testing.T) {
		rawJson := []byte(`
		{
			"type":"regexp",
			"expression":"\\d{1}(\\d{8})\\d+",
			"group_number":1
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
		_, err := NewTransformation(&Params{Expression: `\d(\d+)\d{4}`})
		require.NoError(t, err)

		// invalid regexp sequence **
		_, err = NewTransformation(&Params{Expression: "**"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "error parsing regexp")
	})

	t.Run("Test transformation with one group", func(t *testing.T) {
		tr, err := NewTransformation(&Params{
			Expression:  `number: \d(\d+)\d{4}`,
			GroupNumber: 1,
			Action: action.ActionFunc(func(body []byte) ([]byte, error) {
				require.Equal(t, []byte("23456"), body)
				return body, nil
			}),
		})
		require.NoError(t, err)

		body := []byte("number: 1234567890")
		_, err = tr.Transform(body)
		require.NoError(t, err)
	})

	t.Run("Test transformation with multiple groups", func(t *testing.T) {
		tr, err := NewTransformation(&Params{
			Expression:  `number: (\d+)(\d{4})`,
			GroupNumber: 2,
			Action: action.ActionFunc(func(body []byte) ([]byte, error) {
				require.Equal(t, []byte("7890"), body)
				return body, nil
			}),
		})
		require.NoError(t, err)

		body := []byte("number: 1234567890")
		_, err = tr.Transform(body)
		require.NoError(t, err)
	})

	t.Run("Test transformation when group number exceeds possible number of groups", func(t *testing.T) {
		tr, err := NewTransformation(&Params{
			Expression:  `number: (\d+)(\d{4})`,
			GroupNumber: 5,
			Action: action.ActionFunc(func(body []byte) ([]byte, error) {
				return []byte("xxxx"), nil
			}),
		})
		require.NoError(t, err)

		body := []byte("number: 4242424242424242")
		newBody, err := tr.Transform(body)
		require.NoError(t, err)
		require.Contains(t, string(newBody), "number: 4242424242424242")

		tr2, err := NewTransformation(&Params{
			Expression:  `number: (\d+)(\d{4})`,
			GroupNumber: -1,
			Action: action.ActionFunc(func(body []byte) ([]byte, error) {
				return []byte("xxxx"), nil
			}),
		})
		require.NoError(t, err)

		body = []byte("hello")
		newBody, err = tr2.Transform(body)
		require.NoError(t, err)
		require.Contains(t, string(newBody), "hello")
	})

	t.Run("Test transformation of multiple matches", func(t *testing.T) {
		tr, err := NewTransformation(&Params{
			Expression:  `number: (\d+)(\d{4})`,
			GroupNumber: 2,
			Action: action.ActionFunc(func(body []byte) ([]byte, error) {
				newBody := bytes.Repeat([]byte("x"), len(body))
				return newBody, nil
			}),
		})
		require.NoError(t, err)

		body := []byte("number: 12345 whatever number: 54321")
		want := []byte("number: 1xxxx whatever number: 5xxxx")
		got, err := tr.Transform(body)
		require.NoError(t, err)
		require.Equal(t, string(want), string(got))
	})

	t.Run("Test request transformation", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/request", strings.NewReader(`number: 12345 whatever number: 54321`))
		req.Header.Set("Content-Type", "plain/text")

		tr, err := NewTransformation(&Params{
			Expression:  `number: (\d+)(\d{4})`,
			GroupNumber: 2,
			Action: action.ActionFunc(func(body []byte) ([]byte, error) {
				newBody := bytes.Repeat([]byte("x"), len(body))
				return newBody, nil
			}),
		})
		require.NoError(t, err)

		req, err = tr.TransformRequest(req)
		newBody, err := ioutil.ReadAll(req.Body)
		require.Equal(t, "number: 1xxxx whatever number: 5xxxx", string(newBody))
	})

	t.Run("Test response transformation", func(t *testing.T) {
		res := &http.Response{
			Body:   ioutil.NopCloser(strings.NewReader(`number: 12345 whatever number: 54321`)),
			Header: http.Header{"Content-Type": {"plain/text"}},
		}

		tr, err := NewTransformation(&Params{
			Expression:  `number: (\d+)(\d{4})`,
			GroupNumber: 2,
			Action: action.ActionFunc(func(body []byte) ([]byte, error) {
				newBody := bytes.Repeat([]byte("x"), len(body))
				return newBody, nil
			}),
		})
		require.NoError(t, err)

		res, err = tr.TransformResponse(res)
		require.NoError(t, err)
		newBody, err := ioutil.ReadAll(res.Body)
		require.Equal(t, "number: 1xxxx whatever number: 5xxxx", string(newBody))
	})

}
