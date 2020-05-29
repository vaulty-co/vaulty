package action

import (
	"github.com/rs/xid"
	"github.com/vaulty/vaulty/secrets"
)

type Tokenize struct {
	secretsStorage secrets.SecretsStorage
}

func (a *Tokenize) Transform(body []byte) ([]byte, error) {
	id, _ := xid.New().MarshalText()
	token := append([]byte("tok"), id...)

	err := a.secretsStorage.Set(string(token), body)
	if err != nil {
		return nil, err
	}

	return token, nil
}
