package transform

import (
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
		// I think it may be a good idea to tell something
		// about why we are here. It's because of non-string value or
		// invalid json?
		return body, nil
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
