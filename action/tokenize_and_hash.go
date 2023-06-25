package action

import (
	"fmt"

	"github.com/rs/xid"
	"github.com/vaulty/vaulty/secrets"
	"crypto/sha256"
	"encoding/hex"
)

type TokenizeAndHash struct {
	secretsStorage secrets.Storage
	Format         string
	salt 		   string
}

func (a *TokenizeAndHash) Transform(body []byte) ([]byte, error) {
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

	body = append(body, []byte(a.salt)...)

	sum := sha256.Sum256(body)
	newBody := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(newBody, sum[:])

	hash := string(newBody[:])

	a.secretsStorage.SetWithoutCrypto(hash, token)

	return []byte(token), nil
}

