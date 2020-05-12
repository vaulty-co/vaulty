package transform

import (
	"errors"
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Json struct {
	Action     Transformer
	Expression string
}

func (t *Json) Transform(body []byte) ([]byte, error) {
	result := gjson.GetBytes(body, t.Expression)

	// Currently we perform transformations only over strings
	if result.Type != gjson.String {
		return nil, errors.New(fmt.Sprintf("Result received by expression (%s) is not json String", t.Expression))
	}

	value := result.Str
	newValue, err := t.Action.Transform([]byte(value))
	if err != nil {
		return nil, err
	}

	newBody, err := sjson.SetBytes(body, t.Expression, string(newValue))
	if err != nil {
		return nil, err
	}

	return newBody, nil
}
