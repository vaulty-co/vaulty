package transform

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type Transformer interface {
	Transform(body []byte) ([]byte, error)
}

type TransformerFunc func(body []byte) ([]byte, error)

func (f TransformerFunc) Transform(body []byte) ([]byte, error) {
	return f(body)
}

func Factory(rawInput interface{}, action Transformer) (Transformer, error) {
	input := rawInput.(map[string]interface{})
	switch input["type"] {
	case "json":
		jsonTransformation := &Json{
			Action: action,
		}
		err := mapstructure.Decode(input, jsonTransformation)
		if err != nil {
			return nil, err
		}

		err = jsonTransformation.Validate()
		if err != nil {
			return nil, err
		}

		return jsonTransformation, nil
	case "regexp":
		regexpTransformation := &Regexp{
			Action: action,
		}
		err := mapstructure.Decode(input, regexpTransformation)
		if err != nil {
			return nil, err
		}
		return regexpTransformation, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown transformation type %s", input["type"]))
	}
}
