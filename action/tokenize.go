package action

import (
	"fmt"

	"github.com/rs/xid"
	"github.com/vaulty/vaulty/secrets"
)

type Tokenize struct {
	secretsStorage secrets.SecretsStorage
	Format         string
}

func (a *Tokenize) Transform(body []byte) ([]byte, error) {
	id := xid.New().String()

	var token string

	if a.Format == "email" {
		token = fmt.Sprintf("tok%s@tokenized.local", id)
	} else {
		token = "tok" + id
	}

	err := a.secretsStorage.Set(token, body)
	if err != nil {
		return nil, err
	}

	return []byte(token), nil
}
