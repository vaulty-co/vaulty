package action

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type Action interface {
	Transform(body []byte) ([]byte, error)
}

type ActionFunc func(body []byte) ([]byte, error)

func (a ActionFunc) Transform(body []byte) ([]byte, error) {
	return a(body)
}

func Factory(rawInput interface{}, opts *Options) (Action, error) {
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
			secretsStorage: opts.SecretsStorage,
		}
		err := mapstructure.Decode(input, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	case "detokenize":
		result := &Detokenize{
			secretsStorage: opts.SecretsStorage,
		}
		err := mapstructure.Decode(input, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	case "mask":
		result := &Mask{}
		err := mapstructure.WeakDecode(input, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	case "hash":
		result := &Hash{
			salt: opts.Salt,
		}
		err := mapstructure.Decode(input, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	case "tokenize_and_hash":
		result := &TokenizeAndHash{
			secretsStorage: opts.SecretsStorage,
			salt: opts.Salt,
		}
		err := mapstructure.Decode(input, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown action type %s", input["type"]))
	}
}
