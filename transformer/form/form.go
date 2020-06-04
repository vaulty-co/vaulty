package transform

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/vaulty/vaulty/action"
	"github.com/vaulty/vaulty/transformer"
)

type Transformation struct {
	action action.Action
	fields []string
}

type Params struct {
	Fields string
	Action action.Action
}

var _ transformer.Transformer = (*Transformation)(nil)

func Factory(rawParams map[string]interface{}, act action.Action) (transformer.Transformer, error) {
	params := &Params{
		Action: act,
	}

	err := mapstructure.Decode(rawParams, params)
	if err != nil {
		return nil, err
	}

	return NewTransformation(params)
}

func NewTransformation(params *Params) (*Transformation, error) {
	if params.Fields == "" {
		return nil, errors.New("No fields passed for the form transformation")
	}

	t := &Transformation{
		fields: strings.Split(strings.ReplaceAll(params.Fields, " ", ""), ","),
		action: params.Action,
	}

	return t, nil
}

func (t *Transformation) TransformRequest(req *http.Request) (*http.Request, error) {
	err := t.transformFormData(req)

	return req, err
}

// Currently we do not transform multipart/form-data of the response
func (t *Transformation) TransformResponse(res *http.Response) (*http.Response, error) {
	return res, nil
}

// transformFormData does simple thing. It copies parts
// from the request and writes them into new multipart
// then replaces body of the request
func (t *Transformation) transformFormData(req *http.Request) error {
	// extract boundary parameter from Content-Type header
	header := req.Header.Get("Content-Type")
	_, params, err := mime.ParseMediaType(header)
	if err != nil {
		return err
	}

	boundary, ok := params["boundary"]
	if !ok {
		return fmt.Errorf("boundary was not found in header: %s", header)
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

			newBody, err := t.action.Transform(body)
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
	req.Header.Del("Content-Length")
	req.ContentLength = int64(b.Len())

	return nil
}

func isInSlice(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
