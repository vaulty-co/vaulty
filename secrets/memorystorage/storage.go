package memorystorage

import (
	"github.com/vaulty/vaulty/encryption"
	"github.com/vaulty/vaulty/secrets"
)

var _ secrets.Storage = (*storage)(nil)

type storage struct {
	enc  encryption.Encrypter
	data map[string][]byte
}

func Factory(conf *secrets.Config) (secrets.Storage, error) {
	return New(&Params{
		Encrypter: conf.Encrypter,
	}), nil
}

// Params is used as input to New.
type Params struct {
	// Enc is Encrypter that used to encrypt/decrypt data before putting it
	// into data storage
	Encrypter encryption.Encrypter
}

func New(params *Params) secrets.Storage {
	return &storage{
		data: make(map[string][]byte),
		enc:  params.Encrypter,
	}
}

// Set encrypts valye and adds it to the map. It doesn't check for existing
// value and rewrites if it was set before.
func (s *storage) Set(key string, val []byte) error {
	encrypted, err := s.enc.Encrypt(val)
	if err != nil {
		return err
	}

	s.data[key] = encrypted

	return nil
}

// Set valye and adds it to the map. It doesn't check for existing
// value and rewrites if it was set before.
func (s *storage) SetWithoutCrypto(key string, val string) error {
	s.data[key] = []byte(val)
	return nil
}


// Get decrypts the value from the map and returns it.
func (s *storage) Get(key string) ([]byte, error) {
	encrypted := s.data[key]
	decryted, err := s.enc.Decrypt(encrypted)
	if err != nil {
		return nil, err
	}

	return decryted, nil
}

func (s *storage) Close() error {
	for k := range s.data {
		delete(s.data, k)
	}

	return nil
}
