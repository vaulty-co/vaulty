package secrets

import (
	"io"

	"github.com/vaulty/vaulty/config"
	"github.com/vaulty/vaulty/encryption"
)

type Storage interface {
	Set(key string, val []byte) error
	SetWithoutCrypto(key string, val string) error
	Get(key string) ([]byte, error)

	// Close terminates connections that may remain open. It also
	// may clean up storage memory
	io.Closer
}

type Config struct {
	Encrypter     encryption.Encrypter
	StorageConfig *config.Storage
}

type Factory func(conf *Config) (Storage, error)
