package action

import (
	"github.com/vaulty/vaulty/encrypt"
	"github.com/vaulty/vaulty/secrets"
)

type Options struct {
	Encrypter      encrypt.Encrypter
	SecretsStorage secrets.SecretsStorage
	Salt           string
}
