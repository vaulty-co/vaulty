package json

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/vaulty/vaulty/action"
	"github.com/vaulty/vaulty/transformer"
)

type Transformation struct {
	action      action.Action
	expressions []string
}

type Params struct {
	Expression string
	Action     action.Action
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
	expressions := strings.Split(strings.ReplaceAll(params.Expression, " ", ""), ",")

	for _, exp := range expressions {
		if strings.Count(exp, "#") > 1 {
			return nil, fmt.Errorf("Nested arrays are not supported, but used in the expression: %s", exp)
		}
	}

	t := &Transformation{
		expressions: expressions,
		action:      params.Action,
	}

	return t, nil
}

func (t *Transformation) TransformRequest(req *http.Request) (*http.Request, error) {
	// do nothing if it's not a JSON request
	contentType := req.Header.Get("Content-Type")
	if mediaType, _, _ := mime.ParseMediaType(contentType); mediaType != "application/json" {
		return req, nil
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	newBody, err := t.Transform(body)
	if err != nil {
		return nil, err
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(newBody))
	req.Header.Del("Content-Length")
	req.ContentLength = int64(len(newBody))

	return req, nil
}

func (t *Transformation) TransformResponse(res *http.Response) (*http.Response, error) {
	// do nothing if it's not a JSON request
	contentType := res.Header.Get("Content-Type")
	if mediaType, _, _ := mime.ParseMediaType(contentType); mediaType != "application/json" {
		return res, nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	newBody, err := t.Transform(body)
	if err != nil {
		return nil, err
	}

	res.Body = ioutil.NopCloser(bytes.NewReader(newBody))
	res.Header.Del("Content-Length")
	res.ContentLength = int64(len(newBody))

	return res, nil
}

func (t *Transformation) Transform(body []byte) ([]byte, error) {
	var err error

	for _, expression := range t.expressions {
		body, err = t.TransformWithExpression(body, expression)
		if err != nil {
			return nil, err
		}
	}

	return body, nil
}

func (t *Transformation) TransformWithExpression(body []byte, expression string) ([]byte, error) {
	result := gjson.GetBytes(body, expression)

	// Currently we perform transformations only over strings
	// and arrays (one level only)
	switch {
	case result.Type == gjson.String:
		return t.transformString(body, result, expression)
	case result.IsArray():
		return t.transformArray(body, result, expression)
	default:
		log.Warnf("Unsupported type of json expression result: %s", result.Type)
		return body, nil
	}
}

func (t *Transformation) transformString(body []byte, result gjson.Result, expression string) ([]byte, error) {
	value := result.Str
	newValue, err := t.action.Transform([]byte(value))
	if err != nil {
		return nil, err
	}

	newBody, err := sjson.SetBytes(body, expression, string(newValue))
	if err != nil {
		return nil, err
	}

	return newBody, nil
}

func (t *Transformation) transformArray(body []byte, result gjson.Result, expression string) ([]byte, error) {
	var originalValues []string

	result.ForEach(func(_, res gjson.Result) bool {
		// for non-string values we will add empty string
		// to keep indexes properly
		if res.Type != gjson.String {
			originalValues = append(originalValues, "")
		} else {
			originalValues = append(originalValues, res.String())
		}

		return true
	})

	newBody := body

	for index, value := range originalValues {
		// do not replace empty strings or non-string values
		if value == "" {
			continue
		}

		newValue, err := t.action.Transform([]byte(value))
		if err != nil {
			return nil, err
		}

		setExpression := strings.Replace(expression, "#", strconv.Itoa(index), 1)

		newBody, err = sjson.SetBytes(newBody, setExpression, string(newValue))
		if err != nil {
			return nil, err
		}
	}

	return newBody, nil
}
