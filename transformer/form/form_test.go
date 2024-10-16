package form

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/action"
)

func TestForm(t *testing.T) {
	fakeAction := action.ActionFunc(func(body []byte) ([]byte, error) {
		return append(body, "transformed"...), nil
	})

	t.Run("Test building transformation from JSON", func(t *testing.T) {
		rawJson := []byte(`
		{
			"type":"form",
			"fields":"field1"
		}
		`)

		var input map[string]interface{}
		err := json.Unmarshal(rawJson, &input)

		transformation, err := Factory(input, fakeAction)
		require.NoError(t, err)
		require.NotNil(t, transformation)
	})

	t.Run("Test transformation validation", func(t *testing.T) {
		_, err := NewTransformation(&Params{})
		require.EqualError(t, err, "No fields passed for the form transformation")
	})

	t.Run("Test request transformation of multipart/form-data", func(t *testing.T) {
		formData, err := os.ReadFile("./testdata/form-data.txt")
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

	t.Run("Test request transformation of array of multipart/form-data", func(t *testing.T) {
		formData, err := os.ReadFile("./testdata/form-data.txt")
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/url", bytes.NewReader(formData))
		require.NoError(t, err)
		req.Header = http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}}

		tr, err := NewTransformation(&Params{
			Fields: "field3",
			Action: fakeAction,
		})
		require.NoError(t, err)

		req, err = tr.TransformRequest(req)
		require.NoError(t, err)

		var maxMemory int64 = 32 << 20 // 32MB
		req.ParseMultipartForm(maxMemory)
		require.Equal(t, []string{"value31transformed", "value32transformed"}, req.PostForm["field3"])
	})

	t.Run("Test request transformation of invalid multipart/form-data", func(t *testing.T) {
		formData, err := os.ReadFile("./testdata/invalid-form-data.txt")
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
		require.Error(t, err)
	})

	t.Run("Test request transformation of invalid multipart/form-data", func(t *testing.T) {
		formData, err := os.ReadFile("./testdata/invalid-form-data.txt")
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
		require.Error(t, err)
	})

	t.Run("Test request transformation of application/x-www-form-urlencoded", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/request", strings.NewReader("field1=value1&field2=value2"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		tr, err := NewTransformation(&Params{
			Fields: "field1, field2",
			Action: fakeAction,
		})
		require.NoError(t, err)

		req, err = tr.TransformRequest(req)
		require.NoError(t, err)
		require.Equal(t, "value1transformed", req.FormValue("field1"))
		require.Equal(t, "value2transformed", req.FormValue("field2"))
	})

	t.Run("Test request transformation of array of application/x-www-form-urlencoded", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/request", strings.NewReader("field1=value11&field2=value2&field1=value12"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		tr, err := NewTransformation(&Params{
			Fields: "field1",
			Action: fakeAction,
		})
		require.NoError(t, err)

		req, err = tr.TransformRequest(req)
		require.NoError(t, err)
		req.ParseForm()
		require.Equal(t, []string{"value11transformed", "value12transformed"}, req.Form["field1"])
	})

	t.Run("Test request transformation with unsupported content type", func(t *testing.T) {
		body := io.NopCloser(bytes.NewReader([]byte("{}")))
		req, _ := http.NewRequest("POST", "/request", body)
		req.Header.Set("Content-Type", "application/json")

		tr, err := NewTransformation(&Params{
			Fields: "field1, field2",
			Action: fakeAction,
		})
		require.NoError(t, err)

		newReq, err := tr.TransformRequest(req)
		require.Equal(t, req, newReq)
		require.Equal(t, body, newReq.Body)
	})

	t.Run("Test response transformation does not transform response", func(t *testing.T) {
		res := &http.Response{
			Body: io.NopCloser(bytes.NewReader([]byte("body"))),
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
