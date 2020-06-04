package transform

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestForm(t *testing.T) {
	t.Run("Test building transformer from JSON", func(t *testing.T) {
		postData := `--xxx
Content-Disposition: form-data; name="field1"

value1
--xxx
Content-Disposition: form-data; name="field2"

value2
--xxx
Content-Disposition: form-data; name="file"; filename="file"
Content-Type: application/octet-stream
Content-Transfer-Encoding: binary

binary data
--xxx--
`

		req := &http.Request{
			Method: "POST",
			Header: http.Header{"Content-Type": {`multipart/form-data; boundary=xxx`}},
			Body:   ioutil.NopCloser(strings.NewReader(postData)),
		}

		tr := &Form{
			Fields: "field1, field2",
			Action: TransformerFunc(func(body []byte) ([]byte, error) {
				fmt.Println("Action called")
				return append(body, "transformed"...), nil
			}),
		}

		req, err := tr.Transform(req)
		require.NoError(t, err)

		fmt.Printf("After Value: %s\n", req.PostFormValue("field1"))
		fmt.Printf("After Value: %s\n", req.PostFormValue("field2"))
	})
}
