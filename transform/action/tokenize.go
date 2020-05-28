package action

import (
	"github.com/rs/xid"
	"github.com/vaulty/vaulty/secrets"
)

type Tokenize struct {
	secretStorage secrets.SecretStorage
}

func (a *Tokenize) Transform(body []byte) ([]byte, error) {
	id, _ := xid.New().MarshalText()
	token := append([]byte("tok"), id...)

	err := a.secretStorage.Set(string(token), body)
	if err != nil {
		return nil, err
	}

	return token, nil
}
