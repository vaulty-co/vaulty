package regexp

import (
	"bytes"
	"io"
	"net/http"
	"regexp"

	"github.com/mitchellh/mapstructure"
	"github.com/vaulty/vaulty/action"
	"github.com/vaulty/vaulty/transformer"
)

type Transformation struct {
	expression  string
	groupNumber int
	action      action.Action
	re          *regexp.Regexp
}

type Params struct {
	Expression  string
	GroupNumber int `mapstructure:"group_number"`
	Action      action.Action
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
	var err error

	t := &Transformation{
		action:      params.Action,
		expression:  params.Expression,
		groupNumber: params.GroupNumber,
	}

	t.re, err = regexp.Compile(t.expression)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Transformation) TransformRequest(req *http.Request) (*http.Request, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	newBody, err := t.Transform(body)
	if err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewReader(newBody))
	req.Header.Del("Content-Length")
	req.ContentLength = int64(len(newBody))

	return req, nil
}

func (t *Transformation) TransformResponse(res *http.Response) (*http.Response, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	newBody, err := t.Transform(body)
	if err != nil {
		return nil, err
	}

	res.Body = io.NopCloser(bytes.NewReader(newBody))
	res.Header.Del("Content-Length")
	res.ContentLength = int64(len(newBody))

	return res, nil
}

func (t *Transformation) Transform(body []byte) ([]byte, error) {
	// it does not make sence to do anything
	// if user specified submatch that does not exist
	if t.groupNumber < 1 {
		return body, nil
	}

	results := t.re.FindAllSubmatchIndex(body, -1)

	var newBody []byte

	for _, result := range results {
		// result[2*n:2*n+1] identifies the indexes
		// of the nth submatch.
		// If max position of submatch's end is
		// greater of max position of result it
		// means we don't have enough submatches
		if t.groupNumber*2+1 > len(result)-1 {
			return body, nil
		}

		n := t.groupNumber
		prefix := body[0:result[2*n]]
		value := body[result[2*n]:result[2*n+1]]
		suffix := body[result[2*n+1]:]

		value, err := t.action.Transform(value)
		if err != nil {
			return nil, err
		}

		newBody = append(prefix, value...)
		newBody = append(newBody, suffix...)
	}

	return newBody, nil
}
