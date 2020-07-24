package action

import (
	"github.com/vaulty/vaulty/encryption"
	"github.com/vaulty/vaulty/secrets"
)

type Options struct {
	Encrypter      encryption.Encrypter
	SecretsStorage secrets.Storage
	Salt           string
}
