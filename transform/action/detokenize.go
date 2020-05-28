package action

import (
	"github.com/vaulty/vaulty/secrets"
)

type Detokenize struct {
	secretStorage secrets.SecretStorage
}

func (a *Detokenize) Transform(token []byte) ([]byte, error) {
	val, err := a.secretStorage.Get(string(token))
	if err != nil {
		return nil, err
	}

	return val, nil
}
