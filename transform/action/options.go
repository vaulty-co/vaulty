package action

import (
	"github.com/vaulty/proxy/encrypt"
	"github.com/vaulty/proxy/secrets"
)

type Options struct {
	Encrypter     encrypt.Encrypter
	SecretStorage secrets.SecretStorage
}
