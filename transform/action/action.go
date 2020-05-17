package action

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/vaulty/proxy/transform"
)

func Factory(rawInput interface{}, opts *Options) (transform.Transformer, error) {
	input := rawInput.(map[string]interface{})
	switch input["type"] {
	case "encrypt":
		result := &Encrypt{
			enc: opts.Encrypter,
		}
		err := mapstructure.Decode(input, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	case "decrypt":
		result := &Decrypt{
			enc: opts.Encrypter,
		}
		err := mapstructure.Decode(input, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	case "tokenize":
		result := &Tokenize{
			secretStorage: opts.SecretStorage,
		}
		err := mapstructure.Decode(input, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	case "detokenize":
		result := &Detokenize{
			secretStorage: opts.SecretStorage,
		}
		err := mapstructure.Decode(input, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	case "mask":
		result := &Mask{}
		err := mapstructure.Decode(input, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown action type %s", input["type"]))
	}
}
