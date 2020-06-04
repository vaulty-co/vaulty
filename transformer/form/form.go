package transform

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
)

type Form struct {
	Action Transformer
	Fields string

	fields []string
	once   sync.Once
}

func (t *Form) Transform(req *http.Request) (*http.Request, error) {
	err := t.transformFormData(req)

	return req, err
}

// transformFormData does simple thing. It copies parts
// from the request and writes them into new multipart
// then replaces body of the request
func (t *Form) transformFormData(req *http.Request) error {
	t.once.Do(func() {
		// remove spaces in field if multiple fields are provided
		// for transformation
		t.fields = strings.Split(strings.ReplaceAll(t.Fields, " ", ""), ",")
	})

	// extract boundary parameter from Content-Type header
	v := req.Header.Get("Content-Type")
	_, params, err := mime.ParseMediaType(v)
	if err != nil {
		return err
	}

	boundary, ok := params["boundary"]
	if !ok {
		return fmt.Errorf("boundary was not found in header: %s", v)
	}

	mr := multipart.NewReader(req.Body, boundary)

	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	mw.SetBoundary(boundary)

	for {
		part, err := mr.NextRawPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// create new part
		pw, err := mw.CreatePart(part.Header)

		// if part for the field we have to transform
		if isInSlice(t.fields, part.FormName()) {
			body, err := ioutil.ReadAll(part)
			if err != nil {
				return err
			}

			newBody, err := t.Action.Transform(body)
			if err != nil {
				return err
			}

			io.Copy(pw, bytes.NewBuffer(newBody))
		} else {
			// copy part without modifications
			io.Copy(pw, part)
		}
	}
	mw.Close()

	req.Body = ioutil.NopCloser(bufio.NewReader(&b))
	return nil
}

func isInSlice(slice []string, str string) bool {
	fmt.Printf("Check if %s is in %v\n", str, slice)
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
