package transform

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Json struct {
	Action     Transformer
	Expression string
}

func (t *Json) Validate() error {
	// we do not support nested arrays currently
	if strings.Count(t.Expression, "#") > 1 {
		return fmt.Errorf("Nested arrays are not supported, but used in the expression: %s", t.Expression)
	}

	return nil
}

func (t *Json) Transform(body []byte) ([]byte, error) {
	result := gjson.GetBytes(body, t.Expression)

	// Currently we perform transformations only over strings
	// and arrays (one level only)
	switch {
	case result.Type == gjson.String:
		return t.transformString(body, result)
	case result.IsArray():
		return t.transformArray(body, result)
	default:
		logrus.Warnf("Unsupported type of json expression result: %s", result.Type)
		return body, nil
	}
}

func (t *Json) transformString(body []byte, result gjson.Result) ([]byte, error) {
	value := result.Str
	newValue, err := t.Action.Transform([]byte(value))
	if err != nil {
		return body, nil
	}

	newBody, err := sjson.SetBytes(body, t.Expression, string(newValue))
	if err != nil {
		return nil, err
	}

	return newBody, nil
}

func (t *Json) transformArray(body []byte, result gjson.Result) ([]byte, error) {
	var originalValues []string

	result.ForEach(func(_, res gjson.Result) bool {
		// for non-string values we will add empty string
		// to keep indexes properly
		if res.Type != gjson.String {
			originalValues = append(originalValues, "")
		}

		originalValues = append(originalValues, res.String())

		return true
	})

	newBody := body

	for index, value := range originalValues {
		newValue, err := t.Action.Transform([]byte(value))
		if err != nil {
			return nil, err
		}

		setExpression := strings.Replace(t.Expression, "#", strconv.Itoa(index), 1)

		newBody, err = sjson.SetBytes(newBody, setExpression, string(newValue))
		if err != nil {
			return nil, err
		}
	}

	return newBody, nil
}
