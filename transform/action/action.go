package action

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/vaulty/proxy/transform"
)

func Factory(rawInput interface{}) (transform.Transformer, error) {
	input := rawInput.(map[string]interface{})
	switch input["type"] {
	case "encrypt":
		result := &Encrypt{}
		err := mapstructure.Decode(input, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	case "tokenize":
		result := &Tokenize{}
		err := mapstructure.Decode(input, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown action type %s", input["type"]))
	}
}
