package encryption

import "github.com/vaulty/vaulty/config"

type Encrypter interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

type Factory func(config *config.Config) (Encrypter, error)

var Factories = map[string]Factory{}
