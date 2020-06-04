package transform

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/action"
)

func TestForm(t *testing.T) {
	fakeAction := action.ActionFunc(func(body []byte) ([]byte, error) {
		return append(body, "transformed"...), nil
	})

	t.Run("Test building transformer from JSON", func(t *testing.T) {
		rawJson := []byte(`
		{
			"type":"form",
			"fields":"field1"
		}
		`)

		var input map[string]interface{}
		err := json.Unmarshal(rawJson, &input)

		fakeAction := action.ActionFunc(func(body []byte) ([]byte, error) {
			require.Equal(t, []byte("value1"), body)
			return body, nil
		})

		transformation, err := Factory(input, fakeAction)
		require.NoError(t, err)
		require.NotNil(t, transformation)

		// postData := `--xxx
		// Content-Disposition: form-data; name="field1"

		// value1
		// --xxx
		// Content-Disposition: form-data; name="field2"

		// value2
		// --xxx
		// Content-Disposition: form-data; name="file"; filename="file"
		// Content-Type: application/octet-stream
		// Content-Transfer-Encoding: binary

		// binary data
		// --xxx--
		// `

		// req := httptest.NewRequest("POST", "/url", strings.NewReader(postData))
		// req.Header = http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}}
		// _, err = transformation.TransformRequest(req)
		// require.NoError(t, err)
	})

	t.Run("Test validation", func(t *testing.T) {
		_, err := NewTransformation(&Params{})
		require.EqualError(t, err, "No fields passed for the form transformation")
	})

	t.Run("Test request transformation of multipart/form-data", func(t *testing.T) {
		formData, err := ioutil.ReadFile("./testdata/form-data.txt")
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/url", bytes.NewReader(formData))
		require.NoError(t, err)
		req.Header = http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}}

		tr, err := NewTransformation(&Params{
			Fields: "field1",
			Action: fakeAction,
		})
		require.NoError(t, err)

		req, err = tr.TransformRequest(req)
		require.NoError(t, err)
		require.Equal(t, "value1transformed", req.FormValue("field1"))
	})

	t.Run("Test request transformation of application/x-www-form-urlencoded", func(t *testing.T) {
		// req, err := http.NewRequest("POST", "/url", bytes.NewReader(formData))
		// require.NoError(t, err)
		// req.Header = http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}}

		// tr, err := NewTransformation(&Params{
		// 	Fields: "field1",
		// 	Action: fakeAction,
		// })
		// require.NoError(t, err)

		// req, err = tr.TransformRequest(req)
		// require.NoError(t, err)
		// require.Equal(t, "value1transformed", req.FormValue("field1"))
	})

	t.Run("Test multiple fields transformation", func(t *testing.T) {
	})

	t.Run("Test unsupported content type does nothing", func(t *testing.T) {
	})

	t.Run("Test response transformation does not transform response", func(t *testing.T) {
		res := &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader([]byte("body"))),
		}

		tr, err := NewTransformation(&Params{
			Fields: "field1",
			Action: fakeAction,
		})
		require.NoError(t, err)

		newRes, err := tr.TransformResponse(res)
		require.Equal(t, res, newRes)
		require.Equal(t, res.Body, newRes.Body)
	})
}
