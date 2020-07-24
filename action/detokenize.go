package action

import (
	"github.com/vaulty/vaulty/secrets"
)

type Detokenize struct {
	secretsStorage secrets.Storage
}

func (a *Detokenize) Transform(token []byte) ([]byte, error) {
	val, err := a.secretsStorage.Get(string(token))
	if err != nil {
		return nil, err
	}

	return val, nil
}
